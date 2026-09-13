module fetch_pc_stage(
    input clk, input clrn,
    input [31:0] npc, input [31:0] inst, input wpc,
    output reg [31:0] pc, output [31:0] pc4
);
    always @(posedge clk or negedge clrn) begin
        if (!clrn) pc <= 32'b0;
        else if (wpc) pc <= npc;
    end
    assign pc4 = pc + 32'd4;
endmodule
