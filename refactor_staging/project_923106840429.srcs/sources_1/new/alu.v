module alu(
    input  [31:0] a,
    input  [31:0] b,
    input  [3:0]  alu_ctrl,
    output reg [31:0] result,
    output          zero
);
    always @* begin
        result = 32'b0;
        if      (alu_ctrl == 4'h2) result = a + b;
        else if (alu_ctrl == 4'h6) result = a - b;
        else if (alu_ctrl == 4'h0) result = a & b;
        else if (alu_ctrl == 4'h1) result = a | b;
        else if (alu_ctrl == 4'h3) result = a ^ b;
        else if (alu_ctrl == 4'h4) result = b << a[4:0];
        else if (alu_ctrl == 4'h5) result = b >> a[4:0];
        else if (alu_ctrl == 4'h7) result = $signed(b) >>> a[4:0];
        else if (alu_ctrl == 4'hA) result = ($signed(a) < $signed(b));
        else if (alu_ctrl == 4'hC) result = {b[15:0],16'b0};
    end
    assign zero = (result == 32'b0);
endmodule
