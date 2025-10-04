import request from '../index.ts';

export function getPreInfo() {
  return request<{
    testMode: boolean;
  }>('get', 'info/preInfo', null, 'form', { timeout: 5000 });
}

export function getBaseInfo() {
  return request<DiceBaseInfo>('get', 'info/baseInfo', null, 'form', {
    timeout: 5000,
  });
}

export function getLogFetchAndClear() {
  return request<SysLog[]>('get', 'log/fetchAndClear');
}

export function getHello() {
  return request('get', 'info/hello');
}

export function checkSecurity() {
  return request<{ isOk: boolean }>('get', 'signin/checkSecurity');
}

export function postToolOnebot() {
  return request<{
    ok: boolean;
    ip: string;
    errText: string;
  }>('post', '/utils/onebot');
}

export interface DiceBaseInfo {
  appChannel: string;
  version: string;
  versionSimple: string;
  versionNew: string;
  versionNewNote: string;
  versionCode: number;
  versionNewCode: number;
  memoryAlloc: number;
  memoryUsedSys: number;
  uptime: number;
  OS: string;
  arch: string;
  justForTest: boolean;
  containerMode: boolean;
  extraTitle?: string;
}

type SysLog = {
  level: string;
  ts: number;
  caller: string;
  msg: string;
};
