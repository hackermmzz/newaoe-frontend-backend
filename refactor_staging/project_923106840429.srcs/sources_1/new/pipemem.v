module memory_stage(
    input clk, input clrn,
    input [31:0] exm_alu_result, input [31:0] exm_store_wdata,
    input exm_memwrite, input exm_memread, input exm_memtoreg, input exm_regwrite,
    input [4:0] exm_rd, input [4:0] exm_rt_num,
    input mem_wb_regwrite, input [4:0] mem_wb_rd, input [31:0] mem_wb_wb_data,
    output [31:0] mem_alu_result, output [31:0] mem_read_data,
    output mem_memtoreg, output mem_regwrite, output [4:0] mem_rd
);
    reg [31:0] ram [0:63];
    reg [31:0] write_data;
    integer n;
    wire [5:0] word_index = exm_alu_result[7:2];

    initial begin
        for (n=0; n<64; n=n+1) ram[n] = 32'b0;
        ram[20]=32'hA3; ram[21]=32'h27; ram[22]=32'h79; ram[23]=32'h115;
    end
    always @* begin
        write_data = exm_store_wdata;
        if (mem_wb_regwrite && (mem_wb_rd != 0) && (mem_wb_rd == exm_rt_num))
            write_data = mem_wb_wb_data;
    end
    always @(posedge clk) if (exm_memwrite) ram[word_index] <= write_data;
    assign mem_read_data = ram[word_index];
    assign mem_alu_result = exm_alu_result;
    assign mem_memtoreg = exm_memtoreg;
    assign mem_regwrite = exm_regwrite;
    assign mem_rd = exm_rd;
endmodule
