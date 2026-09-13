module memory_writeback_reg(
    input clk, input clrn,
    input [31:0] mem_alu_result, input [31:0] mem_read_data,
    input [4:0] mem_rd, input mem_memtoreg, input mem_regwrite,
    output reg [31:0] wb_alu_result, output reg [31:0] wb_read_data,
    output reg [4:0] wb_rd, output reg wb_memtoreg, output reg wb_regwrite
);
    always @(posedge clk or negedge clrn) begin
        if (!clrn) begin
            wb_alu_result<=0; wb_read_data<=0; wb_rd<=0; wb_memtoreg<=0; wb_regwrite<=0;
        end else begin
            wb_alu_result<=mem_alu_result; wb_read_data<=mem_read_data; wb_rd<=mem_rd;
            wb_memtoreg<=mem_memtoreg; wb_regwrite<=mem_regwrite;
        end
    end
endmodule
