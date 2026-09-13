module execute_stage(
    input [31:0] ex_rs_data, input [31:0] ex_rt_data, input [31:0] ex_imm_out,
    input [4:0] ex_rs, input [4:0] ex_rt, input [4:0] ex_rd, input [5:0] ex_op,
    input [31:0] ex_pc4, input [5:0] ex_op_funct, input [1:0] ex_aluop,
    input [31:0] ex_mem_alu_result, input [4:0] ex_mem_rd, input ex_mem_regwrite,
    input [31:0] mem_wb_wb_data, input [4:0] mem_wb_rd, input mem_wb_regwrite,
    output [31:0] alu_result, output [31:0] ex_rt_out, output [4:0] ex_dst_reg,
    output alu_zero, output [4:0] ex_rt_num
);
    wire shift_op = (ex_aluop == 2'b10) &&
                    ((ex_op_funct==6'h00)||(ex_op_funct==6'h02)||(ex_op_funct==6'h03));
    wire immediate_op = (ex_aluop==2'b00)||(ex_aluop==2'b11);
    wire link_op = (ex_op==6'h03);
    reg [31:0] op_a, op_b, store_value;
    always @* begin
        op_a = ex_rs_data;
        op_b = immediate_op ? ex_imm_out : ex_rt_data;
        if (shift_op) op_a = {27'b0,ex_imm_out[10:6]};
        if (!shift_op) begin
            if (ex_mem_regwrite && ex_mem_rd!=0 && ex_mem_rd==ex_rs) op_a=ex_mem_alu_result;
            else if (mem_wb_regwrite && mem_wb_rd!=0 && mem_wb_rd==ex_rs) op_a=mem_wb_wb_data;
        end
        if (!immediate_op) begin
            if (ex_mem_regwrite && ex_mem_rd!=0 && ex_mem_rd==ex_rt) op_b=ex_mem_alu_result;
            else if (mem_wb_regwrite && mem_wb_rd!=0 && mem_wb_rd==ex_rt) op_b=mem_wb_wb_data;
        end
        store_value = ex_rt_data;
        if (ex_mem_regwrite && ex_mem_rd!=0 && ex_mem_rd==ex_rt) store_value=ex_mem_alu_result;
        else if (mem_wb_regwrite && mem_wb_rd!=0 && mem_wb_rd==ex_rt) store_value=mem_wb_wb_data;
    end
    wire [5:0] control_field = (ex_aluop==2'b10) ? ex_op_funct : ex_op;
    wire [3:0] control_code;
    wire [31:0] arithmetic_result;
    alu_control decoder(.aluop(ex_aluop),.op_funct(control_field),.alu_ctrl(control_code));
    alu datapath(.a(op_a),.b(op_b),.alu_ctrl(control_code),.result(arithmetic_result),.zero(alu_zero));
    assign alu_result = link_op ? ex_pc4 + 32'd4 : arithmetic_result;
    assign ex_dst_reg = link_op ? 5'd31 : ((ex_aluop==2'b10) ? ex_rd : ex_rt);
    assign ex_rt_out = store_value;
    assign ex_rt_num = ex_rt;
endmodule
