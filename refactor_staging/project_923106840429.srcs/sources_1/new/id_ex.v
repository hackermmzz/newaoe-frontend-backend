module decode_execute_reg(
    input clk, input clrn, input stall, input bubble,
    input [31:0] id_rs_data, input [31:0] id_rt_data, input [31:0] id_imm_out,
    input [4:0] id_rs, input [4:0] id_rt, input [4:0] id_rd, input [5:0] id_op,
    input [31:0] id_pc4,
    input id_regwrite, input id_memread, input id_memwrite, input id_memtoreg,
    input id_aluop1, input id_aluop0,
    output reg [31:0] ex_rs_data, output reg [31:0] ex_rt_data, output reg [31:0] ex_imm_out,
    output reg [4:0] ex_rs, output reg [4:0] ex_rt, output reg [4:0] ex_rd,
    output reg [5:0] ex_op, output reg [31:0] ex_pc4,
    output reg ex_regwrite, output reg ex_memread, output reg ex_memwrite,
    output reg ex_memtoreg, output reg ex_aluop1, output reg ex_aluop0
);
    task automatic zero_stage;
    begin
        ex_rs_data<=0; ex_rt_data<=0; ex_imm_out<=0; ex_rs<=0; ex_rt<=0; ex_rd<=0;
        ex_op<=0; ex_pc4<=0; ex_regwrite<=0; ex_memread<=0; ex_memwrite<=0;
        ex_memtoreg<=0; ex_aluop1<=0; ex_aluop0<=0;
    end endtask
    always @(posedge clk or negedge clrn) begin
        if (!clrn || bubble) zero_stage;
        else if (!stall) begin
            ex_rs_data<=id_rs_data; ex_rt_data<=id_rt_data; ex_imm_out<=id_imm_out;
            ex_rs<=id_rs; ex_rt<=id_rt; ex_rd<=id_rd; ex_op<=id_op; ex_pc4<=id_pc4;
            ex_regwrite<=id_regwrite; ex_memread<=id_memread; ex_memwrite<=id_memwrite;
            ex_memtoreg<=id_memtoreg; ex_aluop1<=id_aluop1; ex_aluop0<=id_aluop0;
        end
    end
endmodule
