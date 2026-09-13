module main_control(
    input  [5:0] op,
    input  [5:0] funct,
    output reg regwrite,
    output reg memread,
    output reg memwrite,
    output reg memtoreg,
    output reg aluop1,
    output reg aluop0,
    output reg jump,
    output reg branch
);
    wire r_format = (op == 6'h00);
    wire jr_code  = r_format && (funct == 6'h08);
    wire addi_code = (op == 6'h08);
    wire logic_imm = (op == 6'h0C) || (op == 6'h0D) ||
                     (op == 6'h0E) || (op == 6'h0F);
    wire load_code  = (op == 6'h23);
    wire store_code = (op == 6'h2B);
    wire branch_code = (op == 6'h04) || (op == 6'h05);
    wire jump_code = (op == 6'h02) || (op == 6'h03);

    always @* begin
        regwrite = ((r_format && !jr_code) || addi_code || logic_imm ||
                    load_code || (op == 6'h03));
        memread  = load_code;
        memwrite = store_code;
        memtoreg = load_code;
        aluop1   = (r_format && !jr_code) || logic_imm;
        aluop0   = branch_code || logic_imm;
        jump     = jump_code || jr_code;
        branch   = branch_code;
    end
endmodule
