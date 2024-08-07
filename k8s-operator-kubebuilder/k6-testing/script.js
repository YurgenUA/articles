import http from 'k6/http';
import { Trend, Gauge } from 'k6/metrics';
import { sleep} from 'k6';

let allocMiBMetric = new Trend('_alloc_MiB', false);
let totalAllocMiBMetric = new Gauge('_totalAlloc_MiB', false);
let sysMiBMetric = new Trend('_sys_MiB', false);
let numGCMetric = new Gauge('_num_GC', false);

export const options = {
  vus: 50,
  duration: '20s',
};

const endpoint_url = 'http://34.116.251.23';

function fetchMemMetrics() {
  let response = http.get(`${endpoint_url}:8080/debug/vars`);
  let debugVarsResponse = JSON.parse(response.body);
  return {
    alloc : debugVarsResponse.memstats.Alloc / (1024 * 1024),
    totalAlloc : debugVarsResponse.memstats.TotalAlloc / (1024 * 1024),
    sys : debugVarsResponse.memstats.Sys / (1024 * 1024),
    numGC: debugVarsResponse.memstats.NumGC
  };
}

export default function() {
  http.get(`${endpoint_url}:8080/hash`);

  let memory = fetchMemMetrics();
  
  allocMiBMetric.add(memory.alloc);
  totalAllocMiBMetric.add(memory.totalAlloc);
  sysMiBMetric.add(memory.sys);
  numGCMetric.add(memory.numGC);

  sleep(1);
}