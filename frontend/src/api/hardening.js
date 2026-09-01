import { apiGet } from "./client.js";

export const getProductionReadiness = () => apiGet("/hardening/readiness");
export const getPerformanceBenchmarks = () => apiGet("/hardening/benchmarks");
