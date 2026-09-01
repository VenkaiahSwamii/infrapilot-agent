import React from 'react';
import { getMachineServices } from '../../../api/machines.js';
import TelemetryTable from './TelemetryTable.jsx';
export default function ServicesTab({ machine }) { return <TelemetryTable machineId={machine?.id} load={getMachineServices} />; }
