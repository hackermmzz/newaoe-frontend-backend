from typing import List, Dict, Any, Optional
import pandas as pd


# ==========================================================
# 自定义行对比函数，自行实现业务逻辑
# ==========================================================
def MergeCMP(row_a: Dict[str, Any], row_b: Dict[str, Any]) -> Dict[str, Any]:
    """
    对比两行同ID数据，返回质量更优的一行数据
    Args:
        row_a: 第一行字典数据，key为列名、value为单元格值
        row_b: 第二行字典数据，key为列名、value为单元格值
    Returns:
        Dict[str, Any]: 对比后更优的单行数据字典
    """
    # TODO: 在此处实现自定义行择优逻辑
    # 示例模板：按分数字段择优
    # return row_a if row_a["score"] >= row_b["score"] else row_b
    if row_a["compile"]==True and row_b["compile"]==False:
        return row_a
    elif row_a["compile"]==False and row_b["compile"]==True:
        return row_b
    if row_a["win"]==True and row_b["win"]==False:
        return row_a
    elif row_a["win"]==False and row_b["win"]==True:
        return row_b
    else:
        return row_a if row_a["score"] >= row_b["score"] else row_b


# ==========================================================
# 多Excel文件合并主逻辑函数
# ==========================================================
def merge_xlsx_by_id(
    files: List[str],
    output: str = "merged.xlsx",
    id_col: str = "id",
    sheet_name: Optional[str] = None,
) -> pd.DataFrame:
    """
    批量合并多个xlsx文件，以指定ID列为主键分组，同ID多行通过MergeCMP择优保留一行
    Args:
        files: 待合并Excel文件路径列表，所有文件列结构必须完全一致
        output: 合并结果输出文件路径；传入None则不写入本地文件
        id_col: 用于分组匹配的主键列名称，默认值id
        sheet_name: 指定读取工作表名称，不传则读取文件第一个sheet
    Returns:
        pd.DataFrame: 合并完成后的完整数据表
    """
    if not files:
        raise ValueError("传入的文件列表files不能为空")

    # 读取所有Excel文件
    df_list = []
    read_config = {"sheet_name": sheet_name} if sheet_name is not None else {}
    for file_path in files:
        df = pd.read_excel(file_path, **read_config)
        if id_col not in df.columns:
            raise ValueError(f"文件 [{file_path}] 中不存在主键列: {id_col}")
        df_list.append(df)

    # 校验所有文件列结构完全统一
    standard_columns = list(df_list[0].columns)
    for idx, df in enumerate(df_list[1:], start=2):
        current_cols = list(df.columns)
        if current_cols != standard_columns:
            raise ValueError(
                f"第 {idx} 个文件与第一个文件列结构不一致\n"
                f"文件1列名: {standard_columns}\n"
                f"文件{idx}列名: {current_cols}"
            )

    # 按ID分组缓存所有行数据，保留ID首次出现顺序
    id_row_buckets: Dict[Any, List[Dict[str, Any]]] = {}
    id_origin_order: List[Any] = []
    for df in df_list:
        for _, row_series in df.iterrows():
            row_dict = row_series.to_dict()
            current_id = row_dict[id_col]
            if current_id not in id_row_buckets:
                id_row_buckets[current_id] = []
                id_origin_order.append(current_id)
            id_row_buckets[current_id].append(row_dict)

    # 每个ID分组迭代择优，生成最终行
    final_row_list: List[Dict[str, Any]] = []
    for target_id in id_origin_order:
        row_group = id_row_buckets[target_id]
        best_row = row_group[0]
        for compare_row in row_group[1:]:
            best_row = MergeCMP(best_row, compare_row)
        final_row_list.append(best_row)

    # 组装结果DataFrame并输出文件
    merged_dataframe = pd.DataFrame(final_row_list, columns=standard_columns)
    if output:
        merged_dataframe.to_excel(output, index=False)

    return merged_dataframe


# ==========================================================
# 命令行运行入口
# ==========================================================
if __name__ == "__main__":
    merge_xlsx_by_id(['C:\\QtHomeWork\\new-aoe-judge\\data2\\result.xlsx',
                      'C:\\QtHomeWork\\new-aoe-judge\\data3\\result.xlsx',
                      'C:\\QtHomeWork\\new-aoe-judge\\data1\\result.xlsx',
                   #   'C:\\QtHomeWork\\new-aoe-judge\\data4\\result.xlsx',
                  #    'C:\\QtHomeWork\\new-aoe-judge\\data5\\result.xlsx',
                   #   "C:\\QtHomeWork\\new-aoe-judge\\merged_1.xlsx"
                      ])
