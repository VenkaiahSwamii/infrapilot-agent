/**
 * Universal helper to derive a normalized, permanent unique identifier for a machine.
 * Correctly distinguishes distinct machines (e.g. Windows host vs WSL2 Linux guest vs VM)
 * even if they share the same base hostname.
 */
export function getMachineId(machine) {
  if (!machine) return '';

  if (typeof machine === 'string' || typeof machine === 'number') {
    return String(machine).toLowerCase().trim();
  }

  // 1. Explicit valid UUID is the highest-fidelity unique machine identity
  const rawId =
    machine.id ||
    machine.ID ||
    machine.Id ||
    machine.machine_id ||
    machine.MachineID ||
    machine.uuid ||
    machine.UUID ||
    '';

  if (rawId && String(rawId).length > 20) {
    return String(rawId).toLowerCase().trim();
  }

  // 2. Derive unique key combining Hostname + OS (e.g. 'venkyyy-windows', 'venkyyy-linux')
  const host = String(machine.hostname || machine.Hostname || machine.name || machine.Name || '').toLowerCase().trim();
  const os = String(machine.os || machine.OS || machine.platform || '').toLowerCase().trim();
  const ip = String(machine.ip_address || machine.IPAddress || '').trim();

  if (host && os) {
    return `${host}-${os}`;
  }

  if (host && ip) {
    return `${host}-${ip}`;
  }

  if (rawId) {
    return String(rawId).toLowerCase().trim();
  }

  return host || 'unknown';
}

/**
 * Compare two machine records to determine if they represent the exact same machine.
 */
export function isSameMachine(a, b) {
  if (!a || !b) return false;
  const idA = a.id || a.ID || a.machine_id;
  const idB = b.id || b.ID || b.machine_id;
  if (idA && idB && String(idA).toLowerCase() === String(idB).toLowerCase()) {
    return true;
  }
  const osA = String(a.os || a.platform || '').toLowerCase();
  const osB = String(b.os || b.platform || '').toLowerCase();
  if (osA && osB && osA !== osB) {
    return false; // Distinct OS is NEVER the same machine
  }
  const ipA = String(a.ip_address || a.IPAddress || '').trim();
  const ipB = String(b.ip_address || b.IPAddress || '').trim();
  if (ipA && ipB && ipA !== ipB) {
    return false; // Distinct IP is NEVER the same machine
  }
  return getMachineId(a) === getMachineId(b);
}
