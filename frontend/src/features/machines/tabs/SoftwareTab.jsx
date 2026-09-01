import React from 'react';
import { getMachineSoftware } from '../../../api/machines.js';
import TelemetryTable from './TelemetryTable.jsx';
export default function SoftwareTab({ machine }) { return <TelemetryTable machineId={machine?.id} load={getMachineSoftware} />; }
