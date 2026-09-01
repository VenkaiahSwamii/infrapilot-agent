import React, { useState, useRef, useEffect } from 'react';

const GRID_COLUMNS = 12;
const ROW_HEIGHT = 80;
const MARGIN = 8;

export default function GridCanvas({ widgets, onUpdateWidget, onDeleteWidget, onAddWidget, onSelectWidget, selectedWidgetId }) {
  const [dragOver, setDragOver] = useState(false);
  const canvasRef = useRef(null);

  const handleDragOver = (e) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'copy';
    setDragOver(true);
  };

  const handleDragLeave = () => {
    setDragOver(false);
  };

  const handleDrop = (e) => {
    e.preventDefault();
    setDragOver(false);

    try {
      const data = JSON.parse(e.dataTransfer.getData('application/json'));
      const rect = canvasRef.current.getBoundingClientRect();
      const x = e.clientX - rect.left;
      const y = e.clientY - rect.top;

      // Calculate grid position
      const cellWidth = (rect.width - MARGIN * (GRID_COLUMNS - 1)) / GRID_COLUMNS;
      const gridX = Math.round(x / (cellWidth + MARGIN));
      const gridY = Math.round(y / (ROW_HEIGHT + MARGIN));

      // Clamp values
      const clampedX = Math.max(0, Math.min(gridX, GRID_COLUMNS - 2));
      const clampedY = Math.max(0, gridY);

      onAddWidget(data.widgetType, data.label, clampedX, clampedY);
    } catch (err) {
      console.error('Drop error:', err);
    }
  };

  const handleWidgetClick = (widget) => {
    onSelectWidget(widget.id);
  };

  const handleResize = (widget, e) => {
    e.stopPropagation();
    const rect = canvasRef.current.getBoundingClientRect();
    const cellWidth = (rect.width - MARGIN * (GRID_COLUMNS - 1)) / GRID_COLUMNS;

    const startX = e.clientX;
    const startWidth = widget.width;
    const startHeight = widget.height;

    const handleMouseMove = (moveEvent) => {
      const deltaX = moveEvent.clientX - startX;
      const deltaY = moveEvent.clientY - rect.top - (widget.y * (ROW_HEIGHT + MARGIN));

      const newWidth = Math.max(2, Math.min(GRID_COLUMNS - widget.x, startWidth + Math.round(deltaX / (cellWidth + MARGIN))));
      const newHeight = Math.max(2, startHeight + Math.round(deltaY / (ROW_HEIGHT + MARGIN)));

      onUpdateWidget(widget.id, { width: newWidth, height: newHeight });
    };

    const handleMouseUp = () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };

    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);
  };

  const handleDragWidget = (widget, e) => {
    e.stopPropagation();
    const rect = canvasRef.current.getBoundingClientRect();
    const cellWidth = (rect.width - MARGIN * (GRID_COLUMNS - 1)) / GRID_COLUMNS;

    const startX = e.clientX;
    const startY = e.clientY;
    const startWidgetX = widget.x;
    const startWidgetY = widget.y;

    const handleMouseMove = (moveEvent) => {
      const deltaX = moveEvent.clientX - startX;
      const deltaY = moveEvent.clientY - startY;

      const newX = Math.max(0, Math.min(GRID_COLUMNS - widget.width, startWidgetX + Math.round(deltaX / (cellWidth + MARGIN))));
      const newY = Math.max(0, startWidgetY + Math.round(deltaY / (ROW_HEIGHT + MARGIN)));

      onUpdateWidget(widget.id, { x: newX, y: newY });
    };

    const handleMouseUp = () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };

    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);
  };

  // Calculate canvas height based on widgets
  const maxY = widgets.reduce((max, widget) => Math.max(max, widget.y + widget.height), 0);
  const canvasHeight = Math.max(600, (maxY + 1) * (ROW_HEIGHT + MARGIN));

  return (
    <div className="grid-canvas-container">
      <div
        ref={canvasRef}
        className={`grid-canvas ${dragOver ? 'drag-over' : ''}`}
        style={{ height: canvasHeight }}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
      >
        {widgets.map((widget) => (
          <div
            key={widget.id}
            className={`grid-widget ${selectedWidgetId === widget.id ? 'selected' : ''}`}
            style={{
              gridColumn: `${widget.x + 1} / span ${widget.width}`,
              gridRow: `${widget.y + 1} / span ${widget.height}`,
            }}
            onClick={() => handleWidgetClick(widget)}
          >
            <div className="widget-content">
              <div className="widget-header">
                <h4>{widget.title}</h4>
                <button 
                  className="delete-widget-btn"
                  onClick={(e) => {
                    e.stopPropagation();
                    onDeleteWidget(widget.id);
                  }}
                >
                  ×
                </button>
              </div>
              <div className="widget-body">
                <WidgetRenderer widget={widget} />
              </div>
            </div>

            <div className="resize-handle" onMouseDown={(e) => handleResize(widget, e)} />
            <div className="drag-handle" onMouseDown={(e) => handleDragWidget(widget, e)} />
          </div>
        ))}

        {widgets.length === 0 && (
          <div className="empty-canvas">
            <p>Drag widgets from the library to start building your dashboard</p>
          </div>
        )}
      </div>
    </div>
  );
}

function WidgetRenderer({ widget }) {
  // Placeholder widget renderer - in production, this would render actual charts
  return (
    <div className="widget-renderer">
      <div className="widget-placeholder">
        <span className="widget-type">{widget.widget_type}</span>
        <p>Widget content will be rendered here</p>
      </div>
    </div>
  );
}