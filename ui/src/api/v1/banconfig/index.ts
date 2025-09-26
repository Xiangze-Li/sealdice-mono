import type { BanConfig } from '#';
import { createRequest } from '../index.ts';
import type { UploadFileInfo } from 'naive-ui';

const baseUrl = '/banconfig/';
const request = createRequest(baseUrl);

export function getBanConfig() {
  return request<BanConfig>('get', 'get');
}

export function setBanConfig(config: BanConfig) {
  return request('post', 'set', config);
}

export function getBanConfigList() {
  return request<BanConfigItem[]>('get', 'list');
}

export function postMapDelOne(item: BanConfigItem) {
  return request('post', 'map_delete_one', item);
}

export function postMapAddOne(ID: string, rank: number, name: string, reason: string) {
  return request('post', 'map_add_one', {
    ID,
    rank,
    name,
    reasons: reason ? [reason] : [],
  });
}

export function importBanConfig(file: UploadFileInfo | Blob) {
  return request<
    | { result: true }
    | {
        result: false;
        err: string;
      }
  >('post', 'import', { file }, 'formdata');
}

type BanConfigItem = {
  ID: string;
  name: string;
  rank: number;
  reason: string;
  rankText?: string;
};
