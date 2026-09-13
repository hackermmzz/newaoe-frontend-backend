module execute_memory_reg(
    input clk, input clrn,
    input [31:0] alu_result, input [31:0] store_rt_data,
    input [4:0] rd, input [4:0] ex_rt_num,
    input memwrite, input memread, input memtoreg, input regwrite,
    output reg [31:0] exm_alu_result, output reg [31:0] exm_store_wdata,
    output reg [4:0] exm_rd, output reg [4:0] exm_rt_num,
    output reg exm_memwrite, output reg exm_memread,
    output reg exm_memtoreg, output reg exm_regwrite
);
    always @(posedge clk or negedge clrn) begin
        if (!clrn) begin
            exm_alu_result<=0; exm_store_wdata<=0; exm_rd<=0; exm_rt_num<=0;
            exm_memwrite<=0; exm_memread<=0; exm_memtoreg<=0; exm_regwrite<=0;
        end else begin
            exm_alu_result<=alu_result; exm_store_wdata<=store_rt_data;
            exm_rd<=rd; exm_rt_num<=ex_rt_num;
            exm_memwrite<=memwrite; exm_memread<=memread;
            exm_memtoreg<=memtoreg; exm_regwrite<=regwrite;
        end
    end
endmodule
