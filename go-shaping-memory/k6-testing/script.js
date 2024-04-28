import http from 'k6/http';
import { Trend, Gauge } from 'k6/metrics';
import { sleep} from 'k6';

let allocMiBMetric = new Trend('_alloc_MiB', false);
let totalAllocMiBMetric = new Gauge('_totalAlloc_MiB', false);
let sysMiBMetric = new Trend('_sys_MiB', false);
let numGCMetric = new Gauge('_num_GC', false);

export const options = {
  vus: 100,
  duration: '30s',
};

function fetchMemMetrics() {
  let response = http.get('http://localhost:8080/debug/vars');
  let debugVarsResponse = JSON.parse(response.body);
  return {
    alloc : debugVarsResponse.memstats.Alloc / (1024 * 1024),
    totalAlloc : debugVarsResponse.memstats.TotalAlloc / (1024 * 1024),
    sys : debugVarsResponse.memstats.Sys / (1024 * 1024),
    numGC: debugVarsResponse.memstats.NumGC
  };
}

export default function() {
  http.get('http://localhost:8080/hash');

  // Fetch memory footprint value once per second
  let memory = fetchMemMetrics();
  
  allocMiBMetric.add(memory.alloc);
  totalAllocMiBMetric.add(memory.totalAlloc);
  sysMiBMetric.add(memory.sys);
  numGCMetric.add(memory.numGC);

  sleep(1);
}