[Setup]
AppName=InfraPilot Enterprise
AppVersion=1.0.0
DefaultDirName={pf}\InfraPilot
DefaultGroupName=InfraPilot
OutputDir=.
OutputBaseFilename=InfraPilotSetup
Compression=lzma
SolidCompression=yes
PrivilegesRequired=admin
ArchitecturesInstallIn64BitMode=x64

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Files]
Source: "..\backend\Dockerfile.backend"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\backend\.env.example"; DestDir: "{app}\backend"; Flags: ignoreversion
Source: "..\frontend\Dockerfile.frontend"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\agent\Dockerfile.agent"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\docker-compose.yml"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\k8s\*"; DestDir: "{app}\k8s"; Flags: ignoreversion recursesubdirs
Source: "..\docs\*"; DestDir: "{app}\docs"; Flags: ignoreversion recursesubdirs
Source: "..\grafana\*"; DestDir: "{app}\grafana"; Flags: ignoreversion recursesubdirs

[Icons]
Name: "{group}\InfraPilot Dashboard"; Filename: "http://localhost"
Name: "{group}\InfraPilot Documentation"; Filename: "{app}\docs\README.md"
Name: "{group}\Uninstall InfraPilot"; Filename: "{uninstallexe}"

[Code]
function InitializeSetup(): Boolean;
begin
  Result := True;
end;