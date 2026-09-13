module alu_control(
    input  [1:0] aluop,
    input  [5:0] op_funct,
    output reg [3:0] alu_ctrl
);
    localparam ADD=4'h2, SUB=4'h6, AND_OP=4'h0, OR_OP=4'h1;
    localparam XOR_OP=4'h3, SLL=4'h4, SRL=4'h5, SRA=4'h7;
    localparam SLT=4'hA, LUI=4'hC;

    always @* begin
        alu_ctrl = AND_OP;
        if (aluop == 2'b00) begin
            alu_ctrl = ADD;
        end else if (aluop == 2'b01) begin
            alu_ctrl = SUB;
        end else if (aluop == 2'b10) begin
            case (op_funct)
                6'h20: alu_ctrl = ADD;
                6'h22: alu_ctrl = SUB;
                6'h24: alu_ctrl = AND_OP;
                6'h25: alu_ctrl = OR_OP;
                6'h26: alu_ctrl = XOR_OP;
                6'h00: alu_ctrl = SLL;
                6'h02: alu_ctrl = SRL;
                6'h03: alu_ctrl = SRA;
                6'h2A: alu_ctrl = SLT;
                default: alu_ctrl = AND_OP;
            endcase
        end else if (aluop == 2'b11) begin
            case (op_funct)
                6'h0C: alu_ctrl = AND_OP;
                6'h0D: alu_ctrl = OR_OP;
                6'h0E: alu_ctrl = XOR_OP;
                6'h0F: alu_ctrl = LUI;
                default: alu_ctrl = OR_OP;
            endcase
        end
    end
endmodule
