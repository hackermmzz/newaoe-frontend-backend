`timescale 1ns / 1ps

module tb_cpu;

    reg clk;
    reg memclock;
    reg clrn;

    //============================================================
    // CPU observation signals
    //============================================================
    wire [31:0] pc;
    wire [31:0] ins;
    wire [31:0] if_id_ins;

    wire        stall;
    wire        bubble;

    wire [31:0] ealu;
    wire [31:0] malu;
    wire [31:0] walu;
    wire [31:0] mem_read_data;

    wire        idex_memread_o;
    wire [4:0]  idex_rt_o;
    wire [4:0]  rs_o;
    wire [4:0]  rt_o;

    wire [31:0] r4_debug_o;
    wire [31:0] r5_debug_o;
    wire [31:0] r8_debug_o;

    wire        take_branch;
    wire [31:0] npc;

    //============================================================
    // CPU instance
    //============================================================
    mips_pipeline_top u_cpu_top (
        .clk            (clk),
        .clrn           (clrn),

        .pc             (pc),
        .ins            (ins),
        .if_id_ins      (if_id_ins),

        .stall          (stall),
        .bubble         (bubble),

        .ealu           (ealu),
        .malu           (malu),
        .walu           (walu),
        .mem_read_data  (mem_read_data),

        .idex_memread_o (idex_memread_o),
        .idex_rt_o      (idex_rt_o),
        .rs_o           (rs_o),
        .rt_o           (rt_o),

        .r4_debug_o     (r4_debug_o),
        .r5_debug_o     (r5_debug_o),
        .r8_debug_o     (r8_debug_o),

        .take_branch_o  (take_branch),
        .npc_o          (npc)
    );

    //============================================================
    // CPU clock: 20 ns period, 50 MHz
    //============================================================
    initial begin
        clk = 1'b0;

        forever begin
            #10 clk = ~clk;
        end
    end

    //============================================================
    // Memory clock
    //
    // The current memory_stage uses clk for RAM access.
    // memclock is retained only for waveform observation.
    //============================================================
    initial begin
        memclock = 1'b0;

        forever begin
            #20 memclock = ~memclock;
        end
    end

    //============================================================
    // Waveform console output
    //============================================================
    initial begin
        $display(
            "Time | PC       | IF_ID_INS | STALL | BUBBLE | EALU     | MALU     | WALU     | R8_DEBUG"
        );

        $display(
            "-----|----------|-----------|-------|--------|----------|----------|----------|----------"
        );

        forever @(negedge clk) begin
            $strobe(
                "%0t | %h | %h | %b | %b | %h | %h | %h | %h",
                $time,
                pc,
                if_id_ins,
                stall,
                bubble,
                ealu,
                malu,
                walu,
                r8_debug_o
            );
        end
    end

    //============================================================
    // Reset and simulation control
    //============================================================
    initial begin
        // Active-low asynchronous reset
        clrn = 1'b0;

        // Keep reset active for one clock cycle
        #20;
        clrn = 1'b1;

        // Run long enough to observe the complete test program
        #2000;

        $display("==============================================");
        $display("Simulation finished");
        $display("Final PC       = %h", pc);
        $display("Final R4       = %h", r4_debug_o);
        $display("Final R5       = %h", r5_debug_o);
        $display("Final R8       = %h", r8_debug_o);
        $display("==============================================");

        $finish;
    end

endmodule