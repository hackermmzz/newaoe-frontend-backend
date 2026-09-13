module pipeline_reg32(
    input clk, input clrn, input en,
    input [31:0] d, output reg [31:0] q
);
    always @(posedge clk or negedge clrn)
        if (!clrn) q <= 32'b0;
        else if (en) q <= d;
endmodule

module fetch_decode_reg(
    input clk, input clrn, input en,
    input [31:0] pc4_in, input [31:0] ins_in,
    output [31:0] pc4_out, output [31:0] ins_out
);
    pipeline_reg32 p_pc(.clk(clk),.clrn(clrn),.en(en),.d(pc4_in),.q(pc4_out));
    pipeline_reg32 p_ir(.clk(clk),.clrn(clrn),.en(en),.d(ins_in),.q(ins_out));
endmodule
