# 排行榜模块

入口：`/home/ranking`；组件：`src/components/Ranking.vue`。

## 请求与响应

GET `${base_url}/getrank?range=0:9`，携带 Cookie。依照请求示例用 URL 查询参数传递分页范围，不发送 GET JSON 请求体。
每页数量由 `RankingRecordPerPage` 配置，默认 10；第 N 页 beg = (N - 1) * pageSize，end = beg + pageSize - 1，两端均包含。

```json
{
  "maxRankInfo": 100,
  "data": [
    {
      "id": "20260001",
      "success": true,
      "score": 95,
      "frame": 12,
      "submittime": "2026-09-12T12:00:00Z",
      "else": "可包含长文本和换行的描述"
    }
  ]
}
```

胜负字段使用 `success` 布尔值。id 为字符串；score 为整数；frame 为非负整数。提交时间键名暂沿用项目现有的 `submittime`，建议返回带时区的 ISO 时间字符串；页面显示浏览器本地时间。

## 展示及排序

表格依次显示 id、胜利/失败、score、frame、提交时间、else。描述默认截断为 100 字符，可展开全部和收起；成功刷新或翻页后重置展开状态。

排序由后端完成，前端不再重新排序，严格按照 `data.data` 数组顺序显示。

不请求或展示总页数，只显示上一页/下一页。下一页返回空数组时禁用下一页；翻页重新请求，不加载全榜。支持空数组、加载失败重试和 15 秒超时，后端自行处理超出范围。
