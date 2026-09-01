import React, { useEffect, useState, useRef, useMemo } from 'react';
import { 
  Folder, 
  File, 
  ArrowUp, 
  ArrowLeft,
  ArrowRight,
  Plus, 
  Trash2, 
  Download, 
  Upload, 
  RefreshCw, 
  HardDrive, 
  X, 
  FolderPlus,
  FilePlus,
  Search,
  CheckCircle2,
  AlertCircle,
  Edit3,
  Eye,
  Copy,
  Check,
  Save,
  LayoutGrid,
  List,
  ChevronRight,
  FileCode,
  FileText,
  FileArchive,
  FileImage,
  Terminal,
  FolderOpen
} from 'lucide-react';
import { apiClient } from '../../../api/client.js';
import { useDashboardStore } from '../../../store/dashboardStore.jsx';

export default function FilesTab({ machine }) {
  const isWindows = !machine?.os || machine.os.toLowerCase().includes('win');
  // Default to C:\Users on Windows or /home on Linux for immediate full write/upload permissions
  const initialRoot = isWindows ? 'C:\\Users' : '/home';
  
  // Navigation & History State
  const [currentPath, setCurrentPath] = useState(initialRoot);
  const [history, setHistory] = useState([initialRoot]);
  const [historyIndex, setHistoryIndex] = useState(0);
  const [isEditingPath, setIsEditingPath] = useState(false);
  const [rawPathInput, setRawPathInput] = useState(initialRoot);

  // Files & Listing State
  const [files, setFiles] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [viewMode, setViewMode] = useState('list'); // 'list' | 'grid'
  const [sortBy, setSortBy] = useState('name'); // 'name' | 'size' | 'mod_time' | 'type'
  const [sortOrder, setSortOrder] = useState('asc'); // 'asc' | 'desc'
  const [selectedPaths, setSelectedPaths] = useState(new Set());

  // Drag and Drop Upload State
  const [isDragging, setIsDragging] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(null);

  // Modals State
  const [showFolderModal, setShowFolderModal] = useState(false);
  const [newFolderName, setNewFolderName] = useState('');
  const [folderError, setFolderError] = useState('');
  const [folderCreating, setFolderCreating] = useState(false);

  const [showFileModal, setShowFileModal] = useState(false);
  const [newFileName, setNewFileName] = useState('');
  const [newFileContent, setNewFileContent] = useState('');
  const [fileModalError, setFileModalError] = useState('');
  const [fileCreating, setFileCreating] = useState(false);

  // Rename Modal
  const [itemToRename, setItemToRename] = useState(null);
  const [renameValue, setRenameValue] = useState('');
  const [renameError, setRenameError] = useState('');
  const [renaming, setRenaming] = useState(false);

  // Delete Confirmation Modal
  const [itemsToDelete, setItemsToDelete] = useState([]); // Array of file objects
  const [deleting, setDeleting] = useState(false);

  // File Viewer & Editor Modal
  const [activeFileViewer, setActiveFileViewer] = useState(null); // { name, path, content, size, isEdited }
  const [viewerLoading, setViewerLoading] = useState(false);
  const [viewerSaving, setViewerSaving] = useState(false);
  const [copiedContent, setCopiedContent] = useState(false);

  const fileInputRef = useRef(null);
  const folderInputRef = useRef(null);
  const newFileInputRef = useRef(null);
  const renameInputRef = useRef(null);
  const dropZoneRef = useRef(null);
  const { addToast } = useDashboardStore();

  // Helper: Normalize path separator based on OS/Path
  const getSep = (path) => (path.includes(':') || path.includes('\\') ? '\\' : '/');

  // Load directory contents
  const loadFiles = async (targetPath, updateHist = true) => {
    if (!machine?.id) return;
    setLoading(true);
    setError('');
    setSelectedPaths(new Set());

    try {
      const response = await apiClient.get(`/machines/${machine.id}/files?path=${encodeURIComponent(targetPath)}`);
      if (Array.isArray(response.data)) {
        setFiles(response.data);
        setCurrentPath(targetPath);
        setRawPathInput(targetPath);

        if (updateHist) {
          setHistory(prev => {
            const next = prev.slice(0, historyIndex + 1);
            if (next[next.length - 1] !== targetPath) {
              next.push(targetPath);
            }
            return next;
          });
          setHistoryIndex(prev => history.length);
        }
      } else {
        setFiles([]);
      }
    } catch (err) {
      console.error('File listing error:', err);
      const msg = err.response?.data?.error || err.message || 'Failed to list directory.';
      setError(msg);
      addToast('critical', 'File Manager Error', msg);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (machine?.id) {
      loadFiles(initialRoot, false);
    }
  }, [machine?.id]);

  // History Navigation
  const canGoBack = historyIndex > 0;
  const canGoForward = historyIndex < history.length - 1;

  const handleGoBack = () => {
    if (canGoBack) {
      const newIdx = historyIndex - 1;
      setHistoryIndex(newIdx);
      loadFiles(history[newIdx], false);
    }
  };

  const handleGoForward = () => {
    if (canGoForward) {
      const newIdx = historyIndex + 1;
      setHistoryIndex(newIdx);
      loadFiles(history[newIdx], false);
    }
  };

  const handleParentDirectory = () => {
    const isWin = currentPath.includes(':') || currentPath.includes('\\');
    const sep = isWin ? '\\' : '/';
    let path = currentPath.replace(isWin ? /\//g : /\\/g, sep);

    if (path.endsWith(sep) && path.length > 1) {
      path = path.slice(0, -1);
    }

    if (isWin && /^[A-Za-z]:\\?$/.test(path)) return;
    if (!isWin && path === '/') return;

    const lastSlash = path.lastIndexOf(sep);
    if (lastSlash === -1) {
      loadFiles(initialRoot);
    } else if (isWin && lastSlash === path.indexOf(':') + 1) {
      loadFiles(path.substring(0, path.indexOf(':') + 1) + '\\');
    } else if (!isWin && lastSlash === 0) {
      loadFiles('/');
    } else {
      loadFiles(path.substring(0, lastSlash) + sep);
    }
  };

  // Breadcrumbs calculation
  const breadcrumbs = useMemo(() => {
    const isWin = currentPath.includes(':') || currentPath.includes('\\');
    const sep = isWin ? '\\' : '/';
    const parts = currentPath.split(/[\\/]+/).filter(Boolean);

    if (isWin) {
      const drive = currentPath.split(':')[0] + ':';
      const trail = [{ name: drive + '\\', path: drive + '\\' }];
      let accumulated = drive + '\\';
      parts.slice(1).forEach(part => {
        accumulated += part + '\\';
        trail.push({ name: part, path: accumulated });
      });
      return trail;
    } else {
      const trail = [{ name: 'root', path: '/' }];
      let accumulated = '';
      parts.forEach(part => {
        accumulated += '/' + part;
        trail.push({ name: part, path: accumulated });
      });
      return trail;
    }
  }, [currentPath]);

  // Sorting & Filtering
  const processedFiles = useMemo(() => {
    let list = [...files];

    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      list = list.filter(f => f.name.toLowerCase().includes(query));
    }

    list.sort((a, b) => {
      // Always directories first
      if (a.type !== b.type) {
        return a.type === 'directory' ? -1 : 1;
      }
      let valA = a[sortBy];
      let valB = b[sortBy];

      if (sortBy === 'size') {
        valA = a.size || 0;
        valB = b.size || 0;
        return sortOrder === 'asc' ? valA - valB : valB - valA;
      }
      if (sortBy === 'mod_time') {
        valA = new Date(a.mod_time || 0).getTime();
        valB = new Date(b.mod_time || 0).getTime();
        return sortOrder === 'asc' ? valA - valB : valB - valA;
      }
      
      // Default name/string sorting
      valA = (a.name || '').toLowerCase();
      valB = (b.name || '').toLowerCase();
      return sortOrder === 'asc' ? valA.localeCompare(valB) : valB.localeCompare(valA);
    });

    return list;
  }, [files, searchQuery, sortBy, sortOrder]);

  const toggleSort = (field) => {
    if (sortBy === field) {
      setSortOrder(prev => (prev === 'asc' ? 'desc' : 'asc'));
    } else {
      setSortBy(field);
      setSortOrder('asc');
    }
  };

  // Selection handlers
  const handleSelectAll = (e) => {
    if (e.target.checked) {
      setSelectedPaths(new Set(processedFiles.map(f => f.path)));
    } else {
      setSelectedPaths(new Set());
    }
  };

  const toggleSelect = (path) => {
    setSelectedPaths(prev => {
      const next = new Set(prev);
      if (next.has(path)) next.delete(path);
      else next.add(path);
      return next;
    });
  };

  // Upload Logic (supports multiple files + drag drop)
  const handleUploadFiles = async (filesList) => {
    if (!filesList || filesList.length === 0) return;
    setUploading(true);
    setUploadProgress(`Uploading ${filesList.length} file(s)...`);

    try {
      for (let i = 0; i < filesList.length; i++) {
        const file = filesList[i];
        setUploadProgress(`Uploading ${file.name} (${i + 1}/${filesList.length})...`);
        const formData = new FormData();
        formData.append('file', file);
        await apiClient.post(`/files/upload?path=${encodeURIComponent(currentPath)}`, formData);
      }
      addToast('success', 'Upload Complete', `Successfully uploaded ${filesList.length} file(s).`);
      loadFiles(currentPath);
    } catch (err) {
      console.error('Upload error:', err);
      addToast('critical', 'Upload Failed', err.response?.data?.error || err.message);
    } finally {
      setUploading(false);
      setUploadProgress(null);
      if (fileInputRef.current) fileInputRef.current.value = '';
    }
  };

  // Drag and Drop listeners
  const handleDragOver = (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (!isDragging) setIsDragging(true);
  };

  const handleDragLeave = (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.relatedTarget === null || !dropZoneRef.current?.contains(e.relatedTarget)) {
      setIsDragging(false);
    }
  };

  const handleDrop = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      handleUploadFiles(e.dataTransfer.files);
    }
  };

  // Create Folder
  const handleCreateFolder = async (e) => {
    if (e) e.preventDefault();
    const cleanName = newFolderName.trim();
    if (!cleanName) {
      setFolderError('Folder name is required.');
      return;
    }
    if (/[<>:"/\\|?*]/.test(cleanName)) {
      setFolderError('Invalid characters in folder name: < > : " / \\ | ? *');
      return;
    }

    const sep = getSep(currentPath);
    const targetPath = (currentPath.endsWith(sep) ? currentPath : currentPath + sep) + cleanName;
    setFolderCreating(true);
    setFolderError('');

    try {
      await apiClient.post('/files/mkdir', { path: targetPath });
      addToast('success', 'Folder Created', `Created ${cleanName}`);
      setShowFolderModal(false);
      setNewFolderName('');
      loadFiles(currentPath);
    } catch (err) {
      setFolderError(err.response?.data?.error || err.message);
    } finally {
      setFolderCreating(false);
    }
  };

  // Create New File
  const handleCreateNewFile = async (e) => {
    if (e) e.preventDefault();
    const cleanName = newFileName.trim();
    if (!cleanName) {
      setFileModalError('File name is required.');
      return;
    }
    if (/[<>:"/\\|?*]/.test(cleanName)) {
      setFileModalError('Invalid characters in file name: < > : " / \\ | ? *');
      return;
    }

    const sep = getSep(currentPath);
    const targetPath = (currentPath.endsWith(sep) ? currentPath : currentPath + sep) + cleanName;
    setFileCreating(true);
    setFileModalError('');

    try {
      await apiClient.post('/files/content', { path: targetPath, content: newFileContent });
      addToast('success', 'File Created', `Created file ${cleanName}`);
      setShowFileModal(false);
      setNewFileName('');
      setNewFileContent('');
      loadFiles(currentPath);
    } catch (err) {
      setFileModalError(err.response?.data?.error || err.message);
    } finally {
      setFileCreating(false);
    }
  };

  // Rename Action
  const handleRenameSubmit = async (e) => {
    if (e) e.preventDefault();
    if (!itemToRename) return;
    const cleanName = renameValue.trim();
    if (!cleanName || cleanName === itemToRename.name) {
      setItemToRename(null);
      return;
    }

    const sep = getSep(itemToRename.path);
    const dir = itemToRename.path.substring(0, itemToRename.path.lastIndexOf(sep));
    const newPath = (dir.endsWith(sep) ? dir : dir + sep) + cleanName;

    setRenaming(true);
    setRenameError('');

    try {
      await apiClient.put('/files/rename', { old_path: itemToRename.path, new_path: newPath });
      addToast('success', 'Renamed Successfully', `${itemToRename.name} → ${cleanName}`);
      setItemToRename(null);
      loadFiles(currentPath);
    } catch (err) {
      setRenameError(err.response?.data?.error || err.message);
    } finally {
      setRenaming(false);
    }
  };

  // Delete Action
  const handleDeleteConfirm = async () => {
    if (itemsToDelete.length === 0) return;
    setDeleting(true);

    try {
      for (const item of itemsToDelete) {
        await apiClient.delete(`/files?path=${encodeURIComponent(item.path)}`);
      }
      addToast('success', 'Deleted Successfully', `Removed ${itemsToDelete.length} item(s).`);
      setItemsToDelete([]);
      loadFiles(currentPath);
    } catch (err) {
      addToast('critical', 'Deletion Error', err.response?.data?.error || err.message);
    } finally {
      setDeleting(false);
    }
  };

  // Download File
  const handleDownload = async (filePath, fileName) => {
    try {
      const response = await apiClient.get(`/files/download?path=${encodeURIComponent(filePath)}`, {
        responseType: 'blob'
      });
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', fileName);
      document.body.appendChild(link);
      link.click();
      link.parentNode.removeChild(link);
      window.URL.revokeObjectURL(url);
      addToast('success', 'Download Started', fileName);
    } catch (err) {
      addToast('critical', 'Download Error', err.response?.data?.error || err.message);
    }
  };

  // View & Edit File
  const handleOpenFileViewer = async (file) => {
    setViewerLoading(true);
    setActiveFileViewer({
      name: file.name,
      path: file.path,
      content: '',
      size: file.size,
      isEdited: false
    });

    try {
      const response = await apiClient.get(`/files/content?path=${encodeURIComponent(file.path)}`);
      setActiveFileViewer(prev => ({
        ...prev,
        content: response.data.content || ''
      }));
    } catch (err) {
      addToast('critical', 'File Read Error', err.response?.data?.error || 'Unable to open file in text viewer');
      setActiveFileViewer(null);
    } finally {
      setViewerLoading(false);
    }
  };

  const handleSaveFileContent = async () => {
    if (!activeFileViewer) return;
    setViewerSaving(true);

    try {
      await apiClient.post('/files/content', {
        path: activeFileViewer.path,
        content: activeFileViewer.content
      });
      addToast('success', 'File Saved', `Changes saved to ${activeFileViewer.name}`);
      setActiveFileViewer(prev => ({ ...prev, isEdited: false }));
      loadFiles(currentPath);
    } catch (err) {
      addToast('critical', 'Save Error', err.response?.data?.error || 'Failed to save changes');
    } finally {
      setViewerSaving(false);
    }
  };

  const handleCopyFileContent = () => {
    if (!activeFileViewer) return;
    navigator.clipboard.writeText(activeFileViewer.content);
    setCopiedContent(true);
    setTimeout(() => setCopiedContent(false), 2000);
  };

  // Helper icon selector
  const getFileIcon = (file) => {
    if (file.type === 'directory') {
      return <Folder size={18} color="#38bdf8" style={{ flexShrink: 0 }} />;
    }
    const ext = file.name.split('.').pop().toLowerCase();
    switch (ext) {
      case 'js':
      case 'jsx':
      case 'ts':
      case 'tsx':
      case 'py':
      case 'go':
      case 'html':
      case 'css':
      case 'json':
      case 'yaml':
      case 'yml':
      case 'sql':
      case 'sh':
      case 'ps1':
      case 'bat':
        return <FileCode size={18} color="#fbbf24" style={{ flexShrink: 0 }} />;
      case 'png':
      case 'jpg':
      case 'jpeg':
      case 'svg':
      case 'gif':
      case 'webp':
      case 'ico':
        return <FileImage size={18} color="#c084fc" style={{ flexShrink: 0 }} />;
      case 'zip':
      case 'tar':
      case 'gz':
      case 'rar':
      case '7z':
        return <FileArchive size={18} color="#f97316" style={{ flexShrink: 0 }} />;
      case 'pdf':
      case 'doc':
      case 'docx':
      case 'txt':
      case 'md':
      case 'log':
      case 'csv':
        return <FileText size={18} color="#34d399" style={{ flexShrink: 0 }} />;
      case 'exe':
      case 'dll':
      case 'bin':
      case 'so':
        return <Terminal size={18} color="#f43f5e" style={{ flexShrink: 0 }} />;
      default:
        return <File size={18} color="#94a3b8" style={{ flexShrink: 0 }} />;
    }
  };

  const formatSize = (bytes) => {
    if (!bytes || bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  // Stats calculation
  const totalFolders = files.filter(f => f.type === 'directory').length;
  const totalFiles = files.filter(f => f.type === 'file').length;
  const totalBytes = files.reduce((acc, f) => acc + (f.size || 0), 0);

  return (
    <div 
      ref={dropZoneRef}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      className="tab-container file-explorer-root" 
      style={{ 
        display: 'flex', 
        flexDirection: 'column', 
        gap: '14px', 
        position: 'relative',
        minHeight: '620px'
      }}
    >
      {/* Drag and Drop Active Overlay */}
      {isDragging && (
        <div style={{
          position: 'absolute',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundColor: 'rgba(6, 182, 212, 0.15)',
          backdropFilter: 'blur(6px)',
          border: '2px dashed #06b6d4',
          borderRadius: '12px',
          zIndex: 50,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          gap: '12px',
          color: '#06b6d4',
          pointerEvents: 'none'
        }}>
          <Upload size={48} className="animate-bounce" />
          <div style={{ fontSize: '18px', fontWeight: 700 }}>Drop files here to upload instantly</div>
          <div style={{ fontSize: '13px', color: 'var(--muted)' }}>Uploading directly to {currentPath}</div>
        </div>
      )}

      {/* Top Header Card */}
      <div style={{
        backgroundColor: '#0d1220',
        border: '1px solid var(--border-soft)',
        borderRadius: '12px',
        padding: '14px 16px',
        display: 'flex',
        flexDirection: 'column',
        gap: '12px'
      }}>
        {/* Navigation Row: History, Quick Drives, Breadcrumbs */}
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px', flexWrap: 'wrap' }}>
          
          {/* Back / Forward History Controls */}
          <div style={{ display: 'flex', gap: '4px' }}>
            <button
              onClick={handleGoBack}
              disabled={!canGoBack || loading}
              style={{
                width: '32px',
                height: '32px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                backgroundColor: 'rgba(255, 255, 255, 0.04)',
                border: '1px solid var(--border-soft)',
                borderRadius: '6px',
                color: canGoBack ? 'var(--text)' : 'var(--muted)',
                cursor: canGoBack ? 'pointer' : 'not-allowed',
                opacity: canGoBack ? 1 : 0.4
              }}
              title="Back"
            >
              <ArrowLeft size={14} />
            </button>
            <button
              onClick={handleGoForward}
              disabled={!canGoForward || loading}
              style={{
                width: '32px',
                height: '32px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                backgroundColor: 'rgba(255, 255, 255, 0.04)',
                border: '1px solid var(--border-soft)',
                borderRadius: '6px',
                color: canGoForward ? 'var(--text)' : 'var(--muted)',
                cursor: canGoForward ? 'pointer' : 'not-allowed',
                opacity: canGoForward ? 1 : 0.4
              }}
              title="Forward"
            >
              <ArrowRight size={14} />
            </button>
            <button
              onClick={handleParentDirectory}
              disabled={loading || currentPath === 'C:\\' || currentPath === 'D:\\' || currentPath === '/'}
              style={{
                width: '32px',
                height: '32px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                backgroundColor: 'rgba(255, 255, 255, 0.04)',
                border: '1px solid var(--border-soft)',
                borderRadius: '6px',
                color: 'var(--text)',
                cursor: 'pointer'
              }}
              title="Up to Parent Directory"
            >
              <ArrowUp size={14} />
            </button>
          </div>

          {/* Quick Drives / Writable User Shortcuts */}
          <div style={{ display: 'flex', gap: '4px', flexWrap: 'wrap' }}>
            {isWindows ? (
              <>
                <button
                  onClick={() => loadFiles('C:\\Users')}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '5px',
                    padding: '0 10px',
                    height: '32px',
                    backgroundColor: currentPath.toUpperCase().startsWith('C:\\USERS') ? 'rgba(6, 182, 212, 0.18)' : 'rgba(255, 255, 255, 0.03)',
                    border: currentPath.toUpperCase().startsWith('C:\\USERS') ? '1px solid #06b6d4' : '1px solid var(--border-soft)',
                    borderRadius: '6px',
                    color: currentPath.toUpperCase().startsWith('C:\\USERS') ? '#06b6d4' : 'var(--muted)',
                    cursor: 'pointer',
                    fontSize: '11px',
                    fontWeight: 700
                  }}
                  title="Users Folder (Full Write/Upload Access)"
                >
                  <FolderOpen size={13} />
                  C:\Users
                </button>
                <button
                  onClick={() => loadFiles('C:\\')}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '5px',
                    padding: '0 10px',
                    height: '32px',
                    backgroundColor: currentPath.toUpperCase() === 'C:\\' || currentPath.toUpperCase() === 'C:' ? 'rgba(6, 182, 212, 0.18)' : 'rgba(255, 255, 255, 0.03)',
                    border: currentPath.toUpperCase() === 'C:\\' || currentPath.toUpperCase() === 'C:' ? '1px solid #06b6d4' : '1px solid var(--border-soft)',
                    borderRadius: '6px',
                    color: currentPath.toUpperCase() === 'C:\\' || currentPath.toUpperCase() === 'C:' ? '#06b6d4' : 'var(--muted)',
                    cursor: 'pointer',
                    fontSize: '11px',
                    fontWeight: 700
                  }}
                >
                  <HardDrive size={13} />
                  C:\
                </button>
                <button
                  onClick={() => loadFiles('D:\\')}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '5px',
                    padding: '0 10px',
                    height: '32px',
                    backgroundColor: currentPath.toUpperCase().startsWith('D:') ? 'rgba(6, 182, 212, 0.18)' : 'rgba(255, 255, 255, 0.03)',
                    border: currentPath.toUpperCase().startsWith('D:') ? '1px solid #06b6d4' : '1px solid var(--border-soft)',
                    borderRadius: '6px',
                    color: currentPath.toUpperCase().startsWith('D:') ? '#06b6d4' : 'var(--muted)',
                    cursor: 'pointer',
                    fontSize: '11px',
                    fontWeight: 700
                  }}
                >
                  <HardDrive size={13} />
                  D:\
                </button>
              </>
            ) : (
              <>
                <button
                  onClick={() => loadFiles('/home')}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '5px',
                    padding: '0 10px',
                    height: '32px',
                    backgroundColor: currentPath.startsWith('/home') ? 'rgba(6, 182, 212, 0.18)' : 'rgba(255, 255, 255, 0.03)',
                    border: currentPath.startsWith('/home') ? '1px solid #06b6d4' : '1px solid var(--border-soft)',
                    borderRadius: '6px',
                    color: currentPath.startsWith('/home') ? '#06b6d4' : 'var(--muted)',
                    cursor: 'pointer',
                    fontSize: '11px',
                    fontWeight: 700
                  }}
                >
                  <FolderOpen size={13} />
                  /home
                </button>
                <button
                  onClick={() => loadFiles('/')}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '5px',
                    padding: '0 10px',
                    height: '32px',
                    backgroundColor: currentPath === '/' ? 'rgba(6, 182, 212, 0.18)' : 'rgba(255, 255, 255, 0.03)',
                    border: currentPath === '/' ? '1px solid #06b6d4' : '1px solid var(--border-soft)',
                    borderRadius: '6px',
                    color: currentPath === '/' ? '#06b6d4' : 'var(--muted)',
                    cursor: 'pointer',
                    fontSize: '11px',
                    fontWeight: 700
                  }}
                >
                  <HardDrive size={13} />
                  / (root)
                </button>
                <button
                  onClick={() => loadFiles('/var/log')}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '5px',
                    padding: '0 10px',
                    height: '32px',
                    backgroundColor: currentPath === '/var/log' ? 'rgba(6, 182, 212, 0.18)' : 'rgba(255, 255, 255, 0.03)',
                    border: currentPath === '/var/log' ? '1px solid #06b6d4' : '1px solid var(--border-soft)',
                    borderRadius: '6px',
                    color: currentPath === '/var/log' ? '#06b6d4' : 'var(--muted)',
                    cursor: 'pointer',
                    fontSize: '11px',
                    fontWeight: 700
                  }}
                >
                  /var/log
                </button>
              </>
            )}
          </div>

          {/* Breadcrumb Trail / Editable Path Bar */}
          <div style={{ flex: 1, minWidth: '220px', position: 'relative' }}>
            {isEditingPath ? (
              <div style={{ display: 'flex', gap: '4px' }}>
                <input
                  type="text"
                  value={rawPathInput}
                  onChange={(e) => setRawPathInput(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                      setIsEditingPath(false);
                      loadFiles(rawPathInput);
                    } else if (e.key === 'Escape') {
                      setIsEditingPath(false);
                      setRawPathInput(currentPath);
                    }
                  }}
                  autoFocus
                  style={{
                    flex: 1,
                    height: '32px',
                    backgroundColor: '#0a0e17',
                    border: '1px solid #06b6d4',
                    borderRadius: '6px',
                    color: '#38bdf8',
                    padding: '0 10px',
                    fontSize: '12px',
                    fontFamily: 'monospace',
                    outline: 'none'
                  }}
                />
                <button
                  onClick={() => {
                    setIsEditingPath(false);
                    loadFiles(rawPathInput);
                  }}
                  style={{
                    padding: '0 10px',
                    height: '32px',
                    backgroundColor: '#06b6d4',
                    border: 'none',
                    borderRadius: '6px',
                    color: '#080c14',
                    fontSize: '11px',
                    fontWeight: 700,
                    cursor: 'pointer'
                  }}
                >
                  Go
                </button>
              </div>
            ) : (
              <div 
                onClick={() => setIsEditingPath(true)}
                style={{
                  height: '32px',
                  backgroundColor: '#0a0e17',
                  border: '1px solid var(--border-soft)',
                  borderRadius: '6px',
                  padding: '0 8px',
                  display: 'flex',
                  alignItems: 'center',
                  gap: '4px',
                  overflowX: 'auto',
                  cursor: 'text',
                  whiteSpace: 'nowrap'
                }}
                title="Click to edit raw path"
              >
                {breadcrumbs.map((crumb, idx) => (
                  <React.Fragment key={crumb.path}>
                    <span
                      onClick={(e) => {
                        e.stopPropagation();
                        loadFiles(crumb.path);
                      }}
                      style={{
                        color: idx === breadcrumbs.length - 1 ? '#38bdf8' : 'var(--muted)',
                        fontWeight: idx === breadcrumbs.length - 1 ? 700 : 500,
                        fontSize: '12px',
                        cursor: 'pointer',
                        padding: '2px 4px',
                        borderRadius: '4px',
                        transition: 'background-color 0.15s'
                      }}
                      onMouseEnter={(e) => e.target.style.backgroundColor = 'rgba(255,255,255,0.06)'}
                      onMouseLeave={(e) => e.target.style.backgroundColor = 'transparent'}
                    >
                      {crumb.name}
                    </span>
                    {idx < breadcrumbs.length - 1 && (
                      <ChevronRight size={12} color="var(--muted)" style={{ opacity: 0.5 }} />
                    )}
                  </React.Fragment>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Action Controls Row: Search, Actions, View Switcher */}
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: '8px', flexWrap: 'wrap' }}>
          
          {/* Search input */}
          <div style={{ position: 'relative', width: '220px' }}>
            <Search size={14} style={{ position: 'absolute', left: '10px', top: '9px', color: 'var(--muted)' }} />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search files..."
              style={{
                width: '100%',
                height: '32px',
                backgroundColor: 'rgba(255, 255, 255, 0.03)',
                border: '1px solid var(--border-soft)',
                borderRadius: '6px',
                color: '#f1f5f9',
                padding: '0 28px 0 30px',
                fontSize: '12px',
                outline: 'none',
                boxSizing: 'border-box'
              }}
            />
            {searchQuery && (
              <button
                onClick={() => setSearchQuery('')}
                style={{
                  position: 'absolute',
                  right: '6px',
                  top: '6px',
                  background: 'none',
                  border: 'none',
                  color: 'var(--muted)',
                  cursor: 'pointer',
                  padding: '2px'
                }}
              >
                <X size={13} />
              </button>
            )}
          </div>

          {/* Right Action Buttons */}
          <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
            
            {/* New Folder Button */}
            <button
              onClick={() => {
                setNewFolderName('');
                setFolderError('');
                setShowFolderModal(true);
              }}
              disabled={loading}
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: '5px',
                padding: '0 12px',
                height: '32px',
                backgroundColor: '#06b6d4',
                border: 'none',
                borderRadius: '6px',
                color: '#080c14',
                cursor: 'pointer',
                fontSize: '12px',
                fontWeight: 700
              }}
            >
              <FolderPlus size={15} />
              New Folder
            </button>

            {/* New File Button */}
            <button
              onClick={() => {
                setNewFileName('');
                setNewFileContent('');
                setFileModalError('');
                setShowFileModal(true);
              }}
              disabled={loading}
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: '5px',
                padding: '0 12px',
                height: '32px',
                backgroundColor: 'rgba(6, 182, 212, 0.12)',
                border: '1px solid rgba(6, 182, 212, 0.4)',
                borderRadius: '6px',
                color: '#06b6d4',
                cursor: 'pointer',
                fontSize: '12px',
                fontWeight: 700
              }}
            >
              <FilePlus size={15} />
              New File
            </button>

            {/* Upload Button */}
            <button
              onClick={() => fileInputRef.current?.click()}
              disabled={loading || uploading}
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: '5px',
                padding: '0 12px',
                height: '32px',
                backgroundColor: 'rgba(255, 255, 255, 0.05)',
                border: '1px solid var(--border-soft)',
                borderRadius: '6px',
                color: 'var(--text)',
                cursor: uploading ? 'not-allowed' : 'pointer',
                fontSize: '12px',
                fontWeight: 600
              }}
            >
              {uploading ? <RefreshCw size={14} className="spin" /> : <Upload size={15} />}
              {uploading ? 'Uploading...' : 'Upload'}
            </button>
            <input
              type="file"
              multiple
              ref={fileInputRef}
              onChange={(e) => handleUploadFiles(e.target.files)}
              style={{ display: 'none' }}
            />

            {/* Refresh Button */}
            <button
              onClick={() => loadFiles(currentPath)}
              disabled={loading}
              style={{
                width: '32px',
                height: '32px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                backgroundColor: 'rgba(255, 255, 255, 0.04)',
                border: '1px solid var(--border-soft)',
                borderRadius: '6px',
                color: 'var(--text)',
                cursor: 'pointer'
              }}
              title="Refresh"
            >
              <RefreshCw size={14} className={loading ? 'spin' : ''} />
            </button>

            {/* View Mode Toggle */}
            <div style={{ display: 'flex', backgroundColor: 'rgba(255, 255, 255, 0.04)', borderRadius: '6px', border: '1px solid var(--border-soft)', padding: '2px' }}>
              <button
                onClick={() => setViewMode('list')}
                style={{
                  padding: '4px 8px',
                  border: 'none',
                  borderRadius: '4px',
                  backgroundColor: viewMode === 'list' ? 'rgba(6, 182, 212, 0.2)' : 'transparent',
                  color: viewMode === 'list' ? '#06b6d4' : 'var(--muted)',
                  cursor: 'pointer',
                  display: 'flex',
                  alignItems: 'center'
                }}
                title="List View"
              >
                <List size={14} />
              </button>
              <button
                onClick={() => setViewMode('grid')}
                style={{
                  padding: '4px 8px',
                  border: 'none',
                  borderRadius: '4px',
                  backgroundColor: viewMode === 'grid' ? 'rgba(6, 182, 212, 0.2)' : 'transparent',
                  color: viewMode === 'grid' ? '#06b6d4' : 'var(--muted)',
                  cursor: 'pointer',
                  display: 'flex',
                  alignItems: 'center'
                }}
                title="Grid View"
              >
                <LayoutGrid size={14} />
              </button>
            </div>

          </div>
        </div>
      </div>

      {/* Uploading Banner */}
      {uploading && (
        <div style={{
          display: 'flex',
          alignItems: 'center',
          gap: '10px',
          padding: '10px 16px',
          backgroundColor: 'rgba(6, 182, 212, 0.1)',
          border: '1px solid rgba(6, 182, 212, 0.3)',
          borderRadius: '8px',
          color: '#06b6d4',
          fontSize: '13px'
        }}>
          <RefreshCw size={16} className="spin" />
          <span>{uploadProgress || 'Uploading files...'}</span>
        </div>
      )}

      {/* Error Banner */}
      {error && (
        <div style={{
          display: 'flex',
          alignItems: 'center',
          gap: '10px',
          padding: '12px 16px',
          backgroundColor: 'rgba(239, 68, 68, 0.1)',
          border: '1px solid rgba(239, 68, 68, 0.3)',
          borderRadius: '8px',
          color: '#f87171',
          fontSize: '13px'
        }}>
          <AlertCircle size={16} />
          <span>{error}</span>
        </div>
      )}

      {/* Multi-Selection Floating Action Bar */}
      {selectedPaths.size > 0 && (
        <div style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          padding: '10px 16px',
          backgroundColor: 'rgba(6, 182, 212, 0.15)',
          border: '1px solid #06b6d4',
          borderRadius: '8px',
          color: '#f1f5f9',
          fontSize: '13px'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
            <CheckCircle2 size={16} color="#06b6d4" />
            <span><strong>{selectedPaths.size}</strong> item(s) selected</span>
          </div>
          <div style={{ display: 'flex', gap: '8px' }}>
            <button
              onClick={() => {
                const selectedItems = files.filter(f => selectedPaths.has(f.path));
                setItemsToDelete(selectedItems);
              }}
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: '5px',
                padding: '4px 12px',
                backgroundColor: '#ef4444',
                border: 'none',
                borderRadius: '6px',
                color: '#fff',
                fontSize: '12px',
                fontWeight: 700,
                cursor: 'pointer'
              }}
            >
              <Trash2 size={13} />
              Delete Selected
            </button>
            <button
              onClick={() => setSelectedPaths(new Set())}
              style={{
                padding: '4px 10px',
                backgroundColor: 'rgba(255, 255, 255, 0.08)',
                border: '1px solid var(--border-soft)',
                borderRadius: '6px',
                color: 'var(--text)',
                fontSize: '12px',
                cursor: 'pointer'
              }}
            >
              Deselect All
            </button>
          </div>
        </div>
      )}

      {/* Main Content Area */}
      {!error && processedFiles.length === 0 && !loading && (
        <div style={{
          padding: '48px 16px',
          textAlign: 'center',
          backgroundColor: '#0d1220',
          border: '1px solid var(--border-soft)',
          borderRadius: '12px',
          color: 'var(--muted)',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          gap: '12px'
        }}>
          <Folder size={40} color="var(--border-soft)" />
          <div style={{ fontSize: '15px', color: '#f1f5f9', fontWeight: 600 }}>
            {searchQuery ? `No files match "${searchQuery}"` : 'This directory is empty'}
          </div>
          <div style={{ fontSize: '13px', maxWidth: '360px' }}>
            Drag and drop files here or click "Upload" or "New Folder" to add items.
          </div>
        </div>
      )}

      {/* List / Table View */}
      {processedFiles.length > 0 && viewMode === 'list' && (
        <div style={{ backgroundColor: '#0d1220', border: '1px solid var(--border-soft)', borderRadius: '12px', overflow: 'hidden' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '13px', textAlign: 'left' }}>
            <thead>
              <tr style={{ borderBottom: '1px solid var(--border-soft)', color: 'var(--muted)', fontSize: '11px', textTransform: 'uppercase', backgroundColor: 'rgba(255, 255, 255, 0.02)' }}>
                <th style={{ width: '36px', padding: '10px 12px', textAlign: 'center' }}>
                  <input
                    type="checkbox"
                    checked={selectedPaths.size > 0 && selectedPaths.size === processedFiles.length}
                    onChange={handleSelectAll}
                    style={{ cursor: 'pointer' }}
                  />
                </th>
                <th 
                  onClick={() => toggleSort('name')}
                  style={{ padding: '10px 14px', cursor: 'pointer', userSelect: 'none' }}
                >
                  Name {sortBy === 'name' && (sortOrder === 'asc' ? '▲' : '▼')}
                </th>
                <th 
                  onClick={() => toggleSort('type')}
                  style={{ width: '100px', padding: '10px 14px', cursor: 'pointer', userSelect: 'none' }}
                >
                  Type {sortBy === 'type' && (sortOrder === 'asc' ? '▲' : '▼')}
                </th>
                <th 
                  onClick={() => toggleSort('size')}
                  style={{ width: '110px', padding: '10px 14px', cursor: 'pointer', userSelect: 'none' }}
                >
                  Size {sortBy === 'size' && (sortOrder === 'asc' ? '▲' : '▼')}
                </th>
                <th 
                  onClick={() => toggleSort('mod_time')}
                  style={{ width: '180px', padding: '10px 14px', cursor: 'pointer', userSelect: 'none' }}
                >
                  Modified {sortBy === 'mod_time' && (sortOrder === 'asc' ? '▲' : '▼')}
                </th>
                <th style={{ width: '160px', padding: '10px 14px', textAlign: 'right' }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {processedFiles.map((file) => {
                const isDir = file.type === 'directory';
                const isSelected = selectedPaths.has(file.path);
                return (
                  <tr
                    key={file.path}
                    style={{
                      borderBottom: '1px solid var(--border-soft)',
                      backgroundColor: isSelected ? 'rgba(6, 182, 212, 0.08)' : 'transparent',
                      transition: 'background-color 0.15s'
                    }}
                    className="file-row"
                  >
                    {/* Checkbox */}
                    <td style={{ padding: '10px 12px', textAlign: 'center' }}>
                      <input
                        type="checkbox"
                        checked={isSelected}
                        onChange={() => toggleSelect(file.path)}
                        style={{ cursor: 'pointer' }}
                      />
                    </td>

                    {/* File Name & Icon */}
                    <td style={{ padding: '10px 14px' }}>
                      <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                        {getFileIcon(file)}
                        {isDir ? (
                          <span
                            onClick={() => loadFiles(file.path)}
                            style={{
                              color: '#38bdf8',
                              cursor: 'pointer',
                              fontWeight: 600,
                              textDecoration: 'none'
                            }}
                            onMouseEnter={(e) => e.target.style.textDecoration = 'underline'}
                            onMouseLeave={(e) => e.target.style.textDecoration = 'none'}
                          >
                            {file.name}
                          </span>
                        ) : (
                          <span 
                            onClick={() => handleOpenFileViewer(file)}
                            style={{ color: '#f1f5f9', cursor: 'pointer' }}
                            onMouseEnter={(e) => e.target.style.color = '#38bdf8'}
                            onMouseLeave={(e) => e.target.style.color = '#f1f5f9'}
                            title="Click to view file content"
                          >
                            {file.name}
                          </span>
                        )}
                      </div>
                    </td>

                    {/* Type */}
                    <td style={{ padding: '10px 14px', color: 'var(--muted)', fontSize: '12px' }}>
                      {isDir ? 'Folder' : (file.name.split('.').pop().toUpperCase() + ' File')}
                    </td>

                    {/* Size */}
                    <td style={{ padding: '10px 14px', color: 'var(--muted)', fontFamily: 'monospace', fontSize: '12px' }}>
                      {isDir ? '-' : formatSize(file.size)}
                    </td>

                    {/* Date Modified */}
                    <td style={{ padding: '10px 14px', color: 'var(--muted)', fontSize: '12px' }}>
                      {file.mod_time ? new Date(file.mod_time).toLocaleString() : '-'}
                    </td>

                    {/* Actions */}
                    <td style={{ padding: '10px 14px', textAlign: 'right' }}>
                      <div style={{ display: 'flex', gap: '6px', justifyContent: 'flex-end', alignItems: 'center' }}>
                        
                        {/* View / Edit Button */}
                        {!isDir && (
                          <button
                            onClick={() => handleOpenFileViewer(file)}
                            style={{
                              background: 'rgba(6, 182, 212, 0.1)',
                              border: '1px solid rgba(6, 182, 212, 0.3)',
                              borderRadius: '4px',
                              color: '#06b6d4',
                              cursor: 'pointer',
                              padding: '4px 7px',
                              display: 'flex',
                              alignItems: 'center',
                              gap: '4px',
                              fontSize: '11px',
                              fontWeight: 600
                            }}
                            title="View / Edit Content"
                          >
                            <Eye size={12} />
                            View
                          </button>
                        )}

                        {/* Download Button */}
                        {!isDir && (
                          <button
                            onClick={() => handleDownload(file.path, file.name)}
                            style={{
                              background: 'rgba(52, 211, 153, 0.1)',
                              border: '1px solid rgba(52, 211, 153, 0.3)',
                              borderRadius: '4px',
                              color: '#34d399',
                              cursor: 'pointer',
                              padding: '4px 7px',
                              display: 'flex',
                              alignItems: 'center',
                              gap: '4px',
                              fontSize: '11px',
                              fontWeight: 600
                            }}
                            title="Download"
                          >
                            <Download size={12} />
                          </button>
                        )}

                        {/* Rename Button */}
                        <button
                          onClick={() => {
                            setItemToRename(file);
                            setRenameValue(file.name);
                            setRenameError('');
                          }}
                          style={{
                            background: 'rgba(255, 255, 255, 0.05)',
                            border: '1px solid var(--border-soft)',
                            borderRadius: '4px',
                            color: 'var(--muted)',
                            cursor: 'pointer',
                            padding: '4px 7px',
                            display: 'flex',
                            alignItems: 'center'
                          }}
                          title="Rename"
                        >
                          <Edit3 size={12} />
                        </button>

                        {/* Delete Button */}
                        <button
                          onClick={() => setItemsToDelete([file])}
                          style={{
                            background: 'rgba(239, 68, 68, 0.1)',
                            border: '1px solid rgba(239, 68, 68, 0.3)',
                            borderRadius: '4px',
                            color: '#f87171',
                            cursor: 'pointer',
                            padding: '4px 7px',
                            display: 'flex',
                            alignItems: 'center'
                          }}
                          title="Delete"
                        >
                          <Trash2 size={12} />
                        </button>

                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}

      {/* Grid View Mode */}
      {processedFiles.length > 0 && viewMode === 'grid' && (
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fill, minmax(160px, 1fr))',
          gap: '12px'
        }}>
          {processedFiles.map((file) => {
            const isDir = file.type === 'directory';
            const isSelected = selectedPaths.has(file.path);
            return (
              <div
                key={file.path}
                onClick={() => {
                  if (isDir) loadFiles(file.path);
                  else handleOpenFileViewer(file);
                }}
                style={{
                  backgroundColor: isSelected ? 'rgba(6, 182, 212, 0.12)' : '#0d1220',
                  border: isSelected ? '1px solid #06b6d4' : '1px solid var(--border-soft)',
                  borderRadius: '10px',
                  padding: '16px 12px',
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  gap: '8px',
                  cursor: 'pointer',
                  position: 'relative',
                  transition: 'all 0.15s'
                }}
                className="grid-file-card"
              >
                {/* Select Checkbox */}
                <input
                  type="checkbox"
                  checked={isSelected}
                  onChange={(e) => {
                    e.stopPropagation();
                    toggleSelect(file.path);
                  }}
                  style={{
                    position: 'absolute',
                    top: '8px',
                    left: '8px',
                    cursor: 'pointer'
                  }}
                />

                <div style={{ transform: 'scale(1.4)', margin: '10px 0' }}>
                  {getFileIcon(file)}
                </div>

                <div style={{
                  fontSize: '12px',
                  fontWeight: 600,
                  color: isDir ? '#38bdf8' : '#f1f5f9',
                  textAlign: 'center',
                  width: '100%',
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  whiteSpace: 'nowrap'
                }}>
                  {file.name}
                </div>

                <div style={{ fontSize: '11px', color: 'var(--muted)', fontFamily: 'monospace' }}>
                  {isDir ? 'Folder' : formatSize(file.size)}
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Bottom Status Bar */}
      <div style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        padding: '10px 16px',
        backgroundColor: '#0a0e17',
        border: '1px solid var(--border-soft)',
        borderRadius: '8px',
        fontSize: '12px',
        color: 'var(--muted)',
        marginTop: 'auto'
      }}>
        <div style={{ display: 'flex', gap: '16px' }}>
          <span><strong>{totalFolders}</strong> folder(s)</span>
          <span><strong>{totalFiles}</strong> file(s)</span>
          <span>Total Size: <strong>{formatSize(totalBytes)}</strong></span>
        </div>
        <div>
          Host: <strong style={{ color: '#06b6d4' }}>{machine?.hostname || machine?.ip_address || 'Connected Node'}</strong>
        </div>
      </div>

      {/* ========================================================================= */}
      {/* MODAL 1: Create New Folder Modal */}
      {/* ========================================================================= */}
      {showFolderModal && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundColor: 'rgba(0, 0, 0, 0.75)',
          backdropFilter: 'blur(4px)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 9999,
          padding: '16px'
        }}>
          <div style={{
            backgroundColor: '#111827',
            border: '1px solid var(--border-soft)',
            borderRadius: '12px',
            width: '100%',
            maxWidth: '440px',
            overflow: 'hidden',
            boxShadow: '0 20px 25px -5px rgba(0,0,0,0.5)'
          }}>
            <div style={{
              padding: '16px 20px',
              borderBottom: '1px solid var(--border-soft)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between'
            }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                <FolderPlus size={18} color="#06b6d4" />
                <h3 style={{ margin: 0, fontSize: '15px', color: '#f1f5f9', fontWeight: 600 }}>Create New Folder</h3>
              </div>
              <button
                onClick={() => setShowFolderModal(false)}
                style={{ background: 'none', border: 'none', color: 'var(--muted)', cursor: 'pointer' }}
              >
                <X size={18} />
              </button>
            </div>
            <form onSubmit={handleCreateFolder} style={{ padding: '20px' }}>
              <div style={{ marginBottom: '12px', fontSize: '12px', color: 'var(--muted)' }}>
                Target: <span style={{ color: '#06b6d4', fontFamily: 'monospace' }}>{currentPath}</span>
              </div>
              <div style={{ marginBottom: '16px' }}>
                <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, color: '#94a3b8', marginBottom: '6px' }}>
                  Folder Name
                </label>
                <input
                  ref={folderInputRef}
                  type="text"
                  value={newFolderName}
                  onChange={(e) => {
                    setNewFolderName(e.target.value);
                    if (folderError) setFolderError('');
                  }}
                  placeholder="e.g. data_archive or logs"
                  autoFocus
                  style={{
                    width: '100%',
                    height: '38px',
                    backgroundColor: '#0a0e17',
                    border: folderError ? '1px solid #ef4444' : '1px solid var(--border-soft)',
                    borderRadius: '8px',
                    color: '#f1f5f9',
                    padding: '0 12px',
                    fontSize: '13px',
                    outline: 'none',
                    boxSizing: 'border-box'
                  }}
                />
                {folderError && (
                  <div style={{ color: '#ef4444', fontSize: '12px', marginTop: '6px' }}>
                    {folderError}
                  </div>
                )}
              </div>
              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '8px' }}>
                <button
                  type="button"
                  onClick={() => setShowFolderModal(false)}
                  disabled={folderCreating}
                  style={{
                    padding: '0 16px',
                    height: '36px',
                    backgroundColor: 'rgba(255, 255, 255, 0.05)',
                    border: '1px solid var(--border-soft)',
                    borderRadius: '8px',
                    color: 'var(--text)',
                    fontSize: '13px',
                    cursor: 'pointer'
                  }}
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={folderCreating || !newFolderName.trim()}
                  style={{
                    padding: '0 18px',
                    height: '36px',
                    backgroundColor: '#06b6d4',
                    border: 'none',
                    borderRadius: '8px',
                    color: '#080c14',
                    fontSize: '13px',
                    fontWeight: 700,
                    cursor: folderCreating || !newFolderName.trim() ? 'not-allowed' : 'pointer',
                    display: 'flex',
                    alignItems: 'center',
                    gap: '6px'
                  }}
                >
                  {folderCreating && <RefreshCw size={14} className="spin" />}
                  {folderCreating ? 'Creating...' : 'Create Folder'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* ========================================================================= */}
      {/* MODAL 2: Create New File Modal */}
      {/* ========================================================================= */}
      {showFileModal && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundColor: 'rgba(0, 0, 0, 0.75)',
          backdropFilter: 'blur(4px)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 9999,
          padding: '16px'
        }}>
          <div style={{
            backgroundColor: '#111827',
            border: '1px solid var(--border-soft)',
            borderRadius: '12px',
            width: '100%',
            maxWidth: '540px',
            overflow: 'hidden',
            boxShadow: '0 20px 25px -5px rgba(0,0,0,0.5)'
          }}>
            <div style={{
              padding: '16px 20px',
              borderBottom: '1px solid var(--border-soft)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between'
            }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                <FilePlus size={18} color="#06b6d4" />
                <h3 style={{ margin: 0, fontSize: '15px', color: '#f1f5f9', fontWeight: 600 }}>Create New File</h3>
              </div>
              <button
                onClick={() => setShowFileModal(false)}
                style={{ background: 'none', border: 'none', color: 'var(--muted)', cursor: 'pointer' }}
              >
                <X size={18} />
              </button>
            </div>
            <form onSubmit={handleCreateNewFile} style={{ padding: '20px' }}>
              <div style={{ marginBottom: '12px', fontSize: '12px', color: 'var(--muted)' }}>
                Target: <span style={{ color: '#06b6d4', fontFamily: 'monospace' }}>{currentPath}</span>
              </div>
              <div style={{ marginBottom: '14px' }}>
                <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, color: '#94a3b8', marginBottom: '6px' }}>
                  File Name
                </label>
                <input
                  ref={newFileInputRef}
                  type="text"
                  value={newFileName}
                  onChange={(e) => {
                    setNewFileName(e.target.value);
                    if (fileModalError) setFileModalError('');
                  }}
                  placeholder="e.g. config.json, script.sh, or notes.txt"
                  autoFocus
                  style={{
                    width: '100%',
                    height: '38px',
                    backgroundColor: '#0a0e17',
                    border: fileModalError ? '1px solid #ef4444' : '1px solid var(--border-soft)',
                    borderRadius: '8px',
                    color: '#f1f5f9',
                    padding: '0 12px',
                    fontSize: '13px',
                    outline: 'none',
                    boxSizing: 'border-box'
                  }}
                />
              </div>
              <div style={{ marginBottom: '16px' }}>
                <label style={{ display: 'block', fontSize: '12px', fontWeight: 600, color: '#94a3b8', marginBottom: '6px' }}>
                  Initial Content (Optional)
                </label>
                <textarea
                  value={newFileContent}
                  onChange={(e) => setNewFileContent(e.target.value)}
                  rows={6}
                  placeholder="# Write your initial file content here..."
                  style={{
                    width: '100%',
                    backgroundColor: '#0a0e17',
                    border: '1px solid var(--border-soft)',
                    borderRadius: '8px',
                    color: '#f1f5f9',
                    padding: '10px 12px',
                    fontSize: '12px',
                    fontFamily: 'monospace',
                    outline: 'none',
                    boxSizing: 'border-box',
                    resize: 'vertical'
                  }}
                />
                {fileModalError && (
                  <div style={{ color: '#ef4444', fontSize: '12px', marginTop: '6px' }}>
                    {fileModalError}
                  </div>
                )}
              </div>
              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '8px' }}>
                <button
                  type="button"
                  onClick={() => setShowFileModal(false)}
                  disabled={fileCreating}
                  style={{
                    padding: '0 16px',
                    height: '36px',
                    backgroundColor: 'rgba(255, 255, 255, 0.05)',
                    border: '1px solid var(--border-soft)',
                    borderRadius: '8px',
                    color: 'var(--text)',
                    fontSize: '13px',
                    cursor: 'pointer'
                  }}
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={fileCreating || !newFileName.trim()}
                  style={{
                    padding: '0 18px',
                    height: '36px',
                    backgroundColor: '#06b6d4',
                    border: 'none',
                    borderRadius: '8px',
                    color: '#080c14',
                    fontSize: '13px',
                    fontWeight: 700,
                    cursor: fileCreating || !newFileName.trim() ? 'not-allowed' : 'pointer',
                    display: 'flex',
                    alignItems: 'center',
                    gap: '6px'
                  }}
                >
                  {fileCreating && <RefreshCw size={14} className="spin" />}
                  {fileCreating ? 'Creating...' : 'Create File'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* ========================================================================= */}
      {/* MODAL 3: Rename Item Modal */}
      {/* ========================================================================= */}
      {itemToRename && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundColor: 'rgba(0, 0, 0, 0.75)',
          backdropFilter: 'blur(4px)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 9999,
          padding: '16px'
        }}>
          <div style={{
            backgroundColor: '#111827',
            border: '1px solid var(--border-soft)',
            borderRadius: '12px',
            width: '100%',
            maxWidth: '420px',
            overflow: 'hidden'
          }}>
            <div style={{ padding: '16px 20px', borderBottom: '1px solid var(--border-soft)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <h3 style={{ margin: 0, fontSize: '15px', color: '#f1f5f9', fontWeight: 600 }}>Rename Item</h3>
              <button onClick={() => setItemToRename(null)} style={{ background: 'none', border: 'none', color: 'var(--muted)', cursor: 'pointer' }}>
                <X size={18} />
              </button>
            </div>
            <form onSubmit={handleRenameSubmit} style={{ padding: '20px' }}>
              <div style={{ marginBottom: '16px' }}>
                <label style={{ display: 'block', fontSize: '12px', color: 'var(--muted)', marginBottom: '6px' }}>
                  New Name
                </label>
                <input
                  ref={renameInputRef}
                  type="text"
                  value={renameValue}
                  onChange={(e) => {
                    setRenameValue(e.target.value);
                    if (renameError) setRenameError('');
                  }}
                  autoFocus
                  style={{
                    width: '100%',
                    height: '38px',
                    backgroundColor: '#0a0e17',
                    border: renameError ? '1px solid #ef4444' : '1px solid var(--border-soft)',
                    borderRadius: '8px',
                    color: '#f1f5f9',
                    padding: '0 12px',
                    fontSize: '13px',
                    outline: 'none',
                    boxSizing: 'border-box'
                  }}
                />
                {renameError && <div style={{ color: '#ef4444', fontSize: '12px', marginTop: '6px' }}>{renameError}</div>}
              </div>
              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '8px' }}>
                <button
                  type="button"
                  onClick={() => setItemToRename(null)}
                  disabled={renaming}
                  style={{
                    padding: '0 16px',
                    height: '36px',
                    backgroundColor: 'rgba(255, 255, 255, 0.05)',
                    border: '1px solid var(--border-soft)',
                    borderRadius: '8px',
                    color: 'var(--text)',
                    fontSize: '13px',
                    cursor: 'pointer'
                  }}
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={renaming || !renameValue.trim()}
                  style={{
                    padding: '0 18px',
                    height: '36px',
                    backgroundColor: '#06b6d4',
                    border: 'none',
                    borderRadius: '8px',
                    color: '#080c14',
                    fontSize: '13px',
                    fontWeight: 700,
                    cursor: 'pointer'
                  }}
                >
                  {renaming ? 'Renaming...' : 'Rename'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* ========================================================================= */}
      {/* MODAL 4: Delete Confirmation Dialog */}
      {/* ========================================================================= */}
      {itemsToDelete.length > 0 && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundColor: 'rgba(0, 0, 0, 0.75)',
          backdropFilter: 'blur(4px)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 9999,
          padding: '16px'
        }}>
          <div style={{
            backgroundColor: '#111827',
            border: '1px solid var(--border-soft)',
            borderRadius: '12px',
            width: '100%',
            maxWidth: '440px',
            overflow: 'hidden'
          }}>
            <div style={{ padding: '20px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '10px', color: '#f87171', marginBottom: '12px' }}>
                <Trash2 size={20} />
                <h3 style={{ margin: 0, fontSize: '16px', fontWeight: 600 }}>Confirm Deletion</h3>
              </div>
              <p style={{ margin: '0 0 16px', fontSize: '13px', color: 'var(--text)', lineHeight: 1.5 }}>
                Are you sure you want to permanently delete{' '}
                {itemsToDelete.length === 1 ? (
                  <strong style={{ color: '#f1f5f9', fontFamily: 'monospace' }}>{itemsToDelete[0].name}</strong>
                ) : (
                  <strong>{itemsToDelete.length} selected items</strong>
                )}
                ? This action cannot be undone.
              </p>
              <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '8px' }}>
                <button
                  type="button"
                  onClick={() => setItemsToDelete([])}
                  disabled={deleting}
                  style={{
                    padding: '0 16px',
                    height: '36px',
                    backgroundColor: 'rgba(255, 255, 255, 0.05)',
                    border: '1px solid var(--border-soft)',
                    borderRadius: '8px',
                    color: 'var(--text)',
                    fontSize: '13px',
                    fontWeight: 600,
                    cursor: 'pointer'
                  }}
                >
                  Cancel
                </button>
                <button
                  type="button"
                  onClick={handleDeleteConfirm}
                  disabled={deleting}
                  style={{
                    padding: '0 18px',
                    height: '36px',
                    backgroundColor: '#ef4444',
                    border: 'none',
                    borderRadius: '8px',
                    color: '#ffffff',
                    fontSize: '13px',
                    fontWeight: 700,
                    cursor: deleting ? 'not-allowed' : 'pointer',
                    display: 'flex',
                    alignItems: 'center',
                    gap: '6px'
                  }}
                >
                  {deleting && <RefreshCw size={14} className="spin" />}
                  {deleting ? 'Deleting...' : 'Delete Permanently'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* ========================================================================= */}
      {/* MODAL 5: Built-in Code/Text Viewer & Editor */}
      {/* ========================================================================= */}
      {activeFileViewer && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundColor: 'rgba(0, 0, 0, 0.85)',
          backdropFilter: 'blur(6px)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 9999,
          padding: '24px'
        }}>
          <div style={{
            backgroundColor: '#0a0e17',
            border: '1px solid var(--border-soft)',
            borderRadius: '12px',
            width: '100%',
            maxWidth: '960px',
            height: '85vh',
            display: 'flex',
            flexDirection: 'column',
            overflow: 'hidden',
            boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.7)'
          }}>
            {/* Viewer Header */}
            <div style={{
              padding: '12px 20px',
              backgroundColor: '#111827',
              borderBottom: '1px solid var(--border-soft)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between'
            }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '10px', minWidth: 0 }}>
                <FileCode size={18} color="#38bdf8" />
                <div style={{ display: 'flex', flexDirection: 'column', minWidth: 0 }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                    <span style={{ fontSize: '14px', fontWeight: 600, color: '#f1f5f9' }}>
                      {activeFileViewer.name}
                    </span>
                    {activeFileViewer.isEdited && (
                      <span style={{ fontSize: '11px', color: '#fbbf24', fontWeight: 600 }}>● Unsaved Changes</span>
                    )}
                  </div>
                  <span style={{ fontSize: '11px', color: 'var(--muted)', fontFamily: 'monospace', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                    {activeFileViewer.path}
                  </span>
                </div>
              </div>

              {/* Viewer Controls */}
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                <button
                  onClick={handleCopyFileContent}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '4px',
                    padding: '0 10px',
                    height: '32px',
                    backgroundColor: 'rgba(255, 255, 255, 0.05)',
                    border: '1px solid var(--border-soft)',
                    borderRadius: '6px',
                    color: copiedContent ? '#34d399' : 'var(--text)',
                    fontSize: '12px',
                    cursor: 'pointer'
                  }}
                  title="Copy content to clipboard"
                >
                  {copiedContent ? <Check size={14} /> : <Copy size={14} />}
                  {copiedContent ? 'Copied' : 'Copy'}
                </button>

                <button
                  onClick={() => handleDownload(activeFileViewer.path, activeFileViewer.name)}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '4px',
                    padding: '0 10px',
                    height: '32px',
                    backgroundColor: 'rgba(255, 255, 255, 0.05)',
                    border: '1px solid var(--border-soft)',
                    borderRadius: '6px',
                    color: 'var(--text)',
                    fontSize: '12px',
                    cursor: 'pointer'
                  }}
                  title="Download File"
                >
                  <Download size={14} />
                  Download
                </button>

                <button
                  onClick={handleSaveFileContent}
                  disabled={viewerSaving || !activeFileViewer.isEdited}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '5px',
                    padding: '0 14px',
                    height: '32px',
                    backgroundColor: activeFileViewer.isEdited ? '#06b6d4' : 'rgba(255, 255, 255, 0.05)',
                    border: activeFileViewer.isEdited ? 'none' : '1px solid var(--border-soft)',
                    borderRadius: '6px',
                    color: activeFileViewer.isEdited ? '#080c14' : 'var(--muted)',
                    fontSize: '12px',
                    fontWeight: 700,
                    cursor: activeFileViewer.isEdited ? 'pointer' : 'default',
                    opacity: activeFileViewer.isEdited ? 1 : 0.5
                  }}
                >
                  {viewerSaving ? <RefreshCw size={13} className="spin" /> : <Save size={13} />}
                  {viewerSaving ? 'Saving...' : 'Save Changes'}
                </button>

                <button
                  onClick={() => setActiveFileViewer(null)}
                  style={{
                    width: '32px',
                    height: '32px',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    backgroundColor: 'rgba(255, 255, 255, 0.05)',
                    border: '1px solid var(--border-soft)',
                    borderRadius: '6px',
                    color: 'var(--muted)',
                    cursor: 'pointer'
                  }}
                >
                  <X size={16} />
                </button>
              </div>
            </div>

            {/* Viewer Body: Code Editor / Text Area */}
            <div style={{ flex: 1, position: 'relative', overflow: 'hidden', display: 'flex' }}>
              {viewerLoading ? (
                <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', width: '100%', color: '#06b6d4', gap: '8px' }}>
                  <RefreshCw size={20} className="spin" />
                  <span>Loading file content...</span>
                </div>
              ) : (
                <textarea
                  value={activeFileViewer.content}
                  onChange={(e) => {
                    setActiveFileViewer(prev => ({
                      ...prev,
                      content: e.target.value,
                      isEdited: true
                    }));
                  }}
                  style={{
                    width: '100%',
                    height: '100%',
                    backgroundColor: '#0a0e17',
                    color: '#e2e8f0',
                    padding: '16px',
                    fontFamily: '"Fira Code", Menlo, Monaco, Consolas, monospace',
                    fontSize: '13px',
                    lineHeight: '1.6',
                    border: 'none',
                    outline: 'none',
                    resize: 'none',
                    boxSizing: 'border-box'
                  }}
                  spellCheck={false}
                />
              )}
            </div>

            {/* Viewer Footer */}
            <div style={{
              padding: '8px 16px',
              backgroundColor: '#111827',
              borderTop: '1px solid var(--border-soft)',
              display: 'flex',
              justifyContent: 'space-between',
              fontSize: '11px',
              color: 'var(--muted)'
            }}>
              <div>Lines: {activeFileViewer.content.split('\n').length} | Characters: {activeFileViewer.content.length}</div>
              <div>Press Escape or Close button to exit editor</div>
            </div>
          </div>
        </div>
      )}

      {/* Internal Custom CSS */}
      <style>{`
        @keyframes spin { to { transform: rotate(360deg); } }
        .spin { animation: spin 0.8s linear infinite; }
        .file-row:hover { background-color: rgba(255, 255, 255, 0.04) !important; }
        .grid-file-card:hover { 
          transform: translateY(-2px); 
          border-color: #06b6d4 !important; 
          box-shadow: 0 8px 20px -6px rgba(0, 0, 0, 0.5);
        }
      `}</style>
    </div>
  );
}