// frontend/student/src/services/v1/reviewApi.js
import request from './request';

export function generateReviewPackage(data) {
  return request({
    url: '/api/v1/student/review/generate', // 注意路径要和 main.go 里注册的一致
    method: 'post',
    data
  });
}

export function getPackageDetail(id) {
  return request({
    url: `/api/v1/student/review/packages/${id}`,
    method: 'get'
  });
}

export function updatePackageItems(id, items) {
  return request({
    url: `/api/v1/student/review/packages/${id}/items`,
    method: 'put',
    data: { items }
  });
}

export function exportReviewPackage(id, format = 'pdf') {
  return request({
    url: `/api/v1/student/review/packages/${id}/export?format=${format}`,
    method: 'get'
  });
}