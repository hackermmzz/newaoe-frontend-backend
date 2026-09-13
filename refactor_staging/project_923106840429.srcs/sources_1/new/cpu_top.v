module mips_pipeline_top(
    input clk, input clrn,
    output [31:0] pc, output [31:0] ins, output [31:0] if_id_ins,
    output stall, output bubble,
    output [31:0] ealu, output [31:0] malu, output [31:0] walu,
    output [31:0] mem_read_data,
    output idex_memread_o, output [4:0] idex_rt_o,
    output [4:0] rs_o, output [4:0] rt_o,
    output [31:0] r4_debug_o, output [31:0] r5_debug_o, output [31:0] r8_debug_o,
    output take_branch_o, output [31:0] npc_o
);
    wire [31:0] pc_q, pc_plus4, instr;
    wire [31:0] if_pc4, if_instr;
    wire [31:0] rs_value, rt_value, immediate, branch_target;
    wire [4:0] rs_idx, rt_idx, rd_idx;
    wire ctl_regwrite, ctl_memread, ctl_memwrite, ctl_memtoreg;
    wire ctl_alu1, ctl_alu0, ctl_jump, ctl_branch, branch_taken;
    wire pipe_stall, pipe_bubble;

    wire [31:0] ex_rs_value, ex_rt_value, ex_imm, ex_pc4;
    wire [4:0] ex_rs, ex_rt, ex_rd;
    wire [5:0] ex_opcode;
    wire ex_regwrite, ex_memread, ex_memwrite, ex_memtoreg, ex_alu1, ex_alu0;
    wire [31:0] alu_value, store_value;
    wire [4:0] alu_dest, store_reg;
    wire alu_zero;
    wire [31:0] xm_alu, xm_store;
    wire [4:0] xm_dest, xm_store_reg;
    wire xm_memwrite, xm_memread, xm_memtoreg, xm_regwrite;
    wire [31:0] mem_alu, mem_value;
    wire [4:0] mem_dest;
    wire mem_memtoreg, mem_regwrite;
    wire [31:0] wb_alu, wb_value, wb_result;
    wire [4:0] wb_dest;
    wire wb_memtoreg, wb_regwrite;

    assign wb_result = wb_memtoreg ? wb_value : wb_alu;
    assign pc = pc_q;
    assign ins = instr;
    assign if_id_ins = if_instr;
    assign stall = pipe_stall;
    assign bubble = pipe_bubble;
    assign ealu = alu_value;
    assign malu = xm_alu;
    assign walu = wb_result;
    assign mem_read_data = mem_value;
    assign idex_memread_o = ex_memread;
    assign idex_rt_o = ex_rt;
    assign rs_o = rs_idx;
    assign rt_o = rt_idx;
    assign r4_debug_o = dbg4;
    assign r5_debug_o = dbg5;
    assign r8_debug_o = dbg8;
    assign take_branch_o = branch_taken;
    assign npc_o = next_pc;

    wire jr = (if_instr[31:26] == 6'b0) && (if_instr[5:0] == 6'h08);
    wire [31:0] jump_target = {pc_q[31:28],if_instr[25:0],2'b0};
    wire [31:0] next_pc = ctl_jump ? (jr ? rs_value : jump_target) :
                                 (branch_taken ? branch_target : pc_plus4);
    wire [31:0] dbg4, dbg5, dbg8;

    instruction_rom imem(.addra(pc_q[7:2]),.douta(instr));
    fetch_pc_stage pc_unit(.clk(clk),.clrn(clrn),.npc(next_pc),.inst(instr),.wpc(~pipe_stall),
                           .pc(pc_q),.pc4(pc_plus4));
    fetch_decode_reg ifid(.clk(clk),.clrn(clrn),.en(~pipe_stall),.pc4_in(pc_plus4),.ins_in(instr),
                          .pc4_out(if_pc4),.ins_out(if_instr));

    decode_stage id_unit(
        .clk(clk),.clrn(clrn),.if_id_ins(if_instr),.if_id_pc4(if_pc4),
        .mem_wb_rd(wb_dest),.wb_data(wb_result),.mem_wb_regwrite(wb_regwrite),
        .ex_alu_result(alu_value),.ex_rd(alu_dest),.ex_regwrite(ex_regwrite),.ex_memread(ex_memread),
        .ex_mem_alu_result(xm_alu),.ex_mem_rd(xm_dest),.ex_mem_regwrite(xm_regwrite),.ex_mem_memread(xm_memread),
        .rs_data(rs_value),.rt_data(rt_value),.imm_out(immediate),.rs(rs_idx),.rt(rt_idx),.rd(rd_idx),
        .regwrite(ctl_regwrite),.memread(ctl_memread),.memwrite(ctl_memwrite),.memtoreg(ctl_memtoreg),
        .aluop1(ctl_alu1),.aluop0(ctl_alu0),.jump(ctl_jump),.branch(ctl_branch),
        .take_branch(branch_taken),.branch_target(branch_target),.r4_debug(dbg4),.r5_debug(dbg5),.r8_debug(dbg8));

    load_use_hazard hazard(.id_rs(rs_idx),.id_rt(rt_idx),.id_ex_rt(ex_rt),.id_ex_memread(ex_memread),
                           .stall(pipe_stall),.bubble(pipe_bubble));
    decode_execute_reg idex(
        .clk(clk),.clrn(clrn),.stall(pipe_stall),.bubble(pipe_bubble),
        .id_rs_data(rs_value),.id_rt_data(rt_value),.id_imm_out(immediate),.id_rs(rs_idx),.id_rt(rt_idx),.id_rd(rd_idx),
        .id_op(if_instr[31:26]),.id_pc4(if_pc4),.id_regwrite(ctl_regwrite),.id_memread(ctl_memread),
        .id_memwrite(ctl_memwrite),.id_memtoreg(ctl_memtoreg),.id_aluop1(ctl_alu1),.id_aluop0(ctl_alu0),
        .ex_rs_data(ex_rs_value),.ex_rt_data(ex_rt_value),.ex_imm_out(ex_imm),.ex_rs(ex_rs),.ex_rt(ex_rt),
        .ex_rd(ex_rd),.ex_op(ex_opcode),.ex_pc4(ex_pc4),.ex_regwrite(ex_regwrite),.ex_memread(ex_memread),
        .ex_memwrite(ex_memwrite),.ex_memtoreg(ex_memtoreg),.ex_aluop1(ex_alu1),.ex_aluop0(ex_alu0));

    execute_stage ex_unit(
        .ex_rs_data(ex_rs_value),.ex_rt_data(ex_rt_value),.ex_imm_out(ex_imm),.ex_rs(ex_rs),.ex_rt(ex_rt),.ex_rd(ex_rd),
        .ex_op(ex_opcode),.ex_pc4(ex_pc4),.ex_op_funct(ex_imm[5:0]),.ex_aluop({ex_alu1,ex_alu0}),
        .ex_mem_alu_result(xm_alu),.ex_mem_rd(xm_dest),.ex_mem_regwrite(xm_regwrite),
        .mem_wb_wb_data(wb_result),.mem_wb_rd(wb_dest),.mem_wb_regwrite(wb_regwrite),
        .alu_result(alu_value),.ex_rt_out(store_value),.ex_dst_reg(alu_dest),.alu_zero(alu_zero),.ex_rt_num(store_reg));
    execute_memory_reg exmem(
        .clk(clk),.clrn(clrn),.alu_result(alu_value),.store_rt_data(store_value),.rd(alu_dest),.ex_rt_num(store_reg),
        .memwrite(ex_memwrite),.memread(ex_memread),.memtoreg(ex_memtoreg),.regwrite(ex_regwrite),
        .exm_alu_result(xm_alu),.exm_store_wdata(xm_store),.exm_rd(xm_dest),.exm_rt_num(xm_store_reg),
        .exm_memwrite(xm_memwrite),.exm_memread(xm_memread),.exm_memtoreg(xm_memtoreg),.exm_regwrite(xm_regwrite));
    memory_stage mem_unit(
        .clk(clk),.clrn(clrn),.exm_alu_result(xm_alu),.exm_store_wdata(xm_store),.exm_memwrite(xm_memwrite),
        .exm_memread(xm_memread),.exm_memtoreg(xm_memtoreg),.exm_regwrite(xm_regwrite),.exm_rd(xm_dest),.exm_rt_num(xm_store_reg),
        .mem_wb_regwrite(wb_regwrite),.mem_wb_rd(wb_dest),.mem_wb_wb_data(wb_result),
        .mem_alu_result(mem_alu),.mem_read_data(mem_value),.mem_memtoreg(mem_memtoreg),.mem_regwrite(mem_regwrite),.mem_rd(mem_dest));
    memory_writeback_reg memwb(
        .clk(clk),.clrn(clrn),.mem_alu_result(mem_alu),.mem_read_data(mem_value),.mem_memtoreg(mem_memtoreg),
        .mem_regwrite(mem_regwrite),.mem_rd(mem_dest),.wb_alu_result(wb_alu),.wb_read_data(wb_value),
        .wb_memtoreg(wb_memtoreg),.wb_regwrite(wb_regwrite),.wb_rd(wb_dest));
endmodule

module instruction_rom(input [5:0] addra, output [31:0] douta);
    reg [31:0] memory [0:63];
    integer i;
    initial begin
        for (i=0; i<64; i=i+1) memory[i]=32'b0;
        memory[0]=32'h3c010000; memory[1]=32'h34240050; memory[2]=32'h0c00001b;
        memory[3]=32'h20050004; memory[4]=32'hac820000; memory[5]=32'h8c890000;
        memory[6]=32'h01244022; memory[7]=32'h20050003; memory[8]=32'h20a5ffff;
        memory[9]=32'h34a8ffff; memory[10]=32'h39085555; memory[11]=32'h2009ffff;
        memory[12]=32'h312affff; memory[13]=32'h01493025; memory[14]=32'h01494026;
        memory[15]=32'h01463824; memory[16]=32'h10a00003; memory[17]=32'h00000000;
        memory[18]=32'h08000008; memory[19]=32'h00000000; memory[20]=32'h2005ffff;
        memory[21]=32'h000543c0; memory[22]=32'h00084400; memory[23]=32'h00084403;
        memory[24]=32'h000843c2; memory[25]=32'h08000019; memory[26]=32'h00000000;
        memory[27]=32'h00004020; memory[28]=32'h8c890000; memory[29]=32'h01094020;
        memory[30]=32'h20a5ffff; memory[31]=32'h14a0fffc; memory[32]=32'h20840004;
        memory[33]=32'h03e00008; memory[34]=32'h00081000;
    end
    assign douta = memory[addra];
endmodule
