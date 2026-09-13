module decode_stage(
    input clk, input clrn,
    input [31:0] if_id_ins, input [31:0] if_id_pc4,
    input [4:0] mem_wb_rd, input [31:0] wb_data, input mem_wb_regwrite,
    input [31:0] ex_alu_result, input [4:0] ex_rd, input ex_regwrite, input ex_memread,
    input [31:0] ex_mem_alu_result, input [4:0] ex_mem_rd,
    input ex_mem_regwrite, input ex_mem_memread,
    output [31:0] rs_data, output [31:0] rt_data, output [31:0] imm_out,
    output [4:0] rs, output [4:0] rt, output [4:0] rd,
    output regwrite, output memread, output memwrite, output memtoreg,
    output aluop1, output aluop0, output jump, output branch,
    output take_branch, output [31:0] branch_target,
    output [31:0] r4_debug, output [31:0] r5_debug, output [31:0] r8_debug
);
    wire [5:0] opcode = if_id_ins[31:26];
    wire [5:0] funct  = if_id_ins[5:0];
    wire [15:0] imm16 = if_id_ins[15:0];
    assign rs = if_id_ins[25:21];
    assign rt = if_id_ins[20:16];
    assign rd = if_id_ins[15:11];
    wire zero_extend = (opcode==6'h0C)||(opcode==6'h0D)||(opcode==6'h0E)||(opcode==6'h0F);
    assign imm_out = zero_extend ? {16'b0,imm16} : {{16{imm16[15]}},imm16};
    assign branch_target = if_id_pc4 + {imm_out[29:0],2'b0};

    wire [31:0] raw_rs, raw_rt;
    regfile rf(.clk(clk),.clrn(clrn),.we(mem_wb_regwrite),.waddr(mem_wb_rd),.wdata(wb_data),
                .raddr1(rs),.raddr2(rt),.rdata1(raw_rs),.rdata2(raw_rt),
                .r4_debug(r4_debug),.r5_debug(r5_debug),.r8_debug(r8_debug));
    assign rs_data = (mem_wb_regwrite && mem_wb_rd!=0 && mem_wb_rd==rs) ? wb_data : raw_rs;
    assign rt_data = (mem_wb_regwrite && mem_wb_rd!=0 && mem_wb_rd==rt) ? wb_data : raw_rt;

    main_control ctl(.op(opcode),.funct(funct),.regwrite(regwrite),.memread(memread),
        .memwrite(memwrite),.memtoreg(memtoreg),.aluop1(aluop1),.aluop0(aluop0),
        .jump(jump),.branch(branch));

    reg [31:0] br_a, br_b;
    always @* begin
        br_a = rs_data; br_b = rt_data;
        if (ex_regwrite && !ex_memread && ex_rd!=0 && ex_rd==rs) br_a=ex_alu_result;
        else if (ex_mem_regwrite && !ex_mem_memread && ex_mem_rd!=0 && ex_mem_rd==rs) br_a=ex_mem_alu_result;
        else if (mem_wb_regwrite && mem_wb_rd!=0 && mem_wb_rd==rs) br_a=wb_data;
        if (ex_regwrite && !ex_memread && ex_rd!=0 && ex_rd==rt) br_b=ex_alu_result;
        else if (ex_mem_regwrite && !ex_mem_memread && ex_mem_rd!=0 && ex_mem_rd==rt) br_b=ex_mem_alu_result;
        else if (mem_wb_regwrite && mem_wb_rd!=0 && mem_wb_rd==rt) br_b=wb_data;
    end
    assign take_branch = branch && ((opcode==6'h05) ? (br_a != br_b) : (br_a == br_b));
endmodule
