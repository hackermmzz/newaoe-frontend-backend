module regfile(
    input clk, input clrn, input we,
    input [4:0] waddr, input [31:0] wdata,
    input [4:0] raddr1, input [4:0] raddr2,
    output [31:0] rdata1, output [31:0] rdata2,
    output [31:0] r4_debug, output [31:0] r5_debug, output [31:0] r8_debug
);
    reg [31:0] bank [0:31];
    integer k;
    always @(posedge clk or negedge clrn) begin
        if (!clrn) begin
            for (k=0; k<32; k=k+1) bank[k] <= 32'b0;
        end else if (we && (waddr != 5'd0)) begin
            bank[waddr] <= wdata;
        end
    end
    assign rdata1 = (raddr1 == 0) ? 32'b0 : bank[raddr1];
    assign rdata2 = (raddr2 == 0) ? 32'b0 : bank[raddr2];
    assign r4_debug = bank[4];
    assign r5_debug = bank[5];
    assign r8_debug = bank[8];
endmodule
