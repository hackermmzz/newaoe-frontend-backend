// Copyright 1986-2018 Xilinx, Inc. All Rights Reserved.
// --------------------------------------------------------------------------------
// Tool Version: Vivado v.2018.3 (win64) Build 2405991 Thu Dec  6 23:38:27 MST 2018
// Date        : Wed Sep  2 16:57:37 2026
// Host        : fuyou38968 running 64-bit major release  (build 9200)
// Command     : write_verilog -force -mode synth_stub
//               D:/vivadoProject/project_923106840535/project_923106840535.srcs/sources_1/ip/data_ram/data_ram_stub.v
// Design      : data_ram
// Purpose     : Stub declaration of top-level module interface
// Device      : xc7k325tffg900-2
// --------------------------------------------------------------------------------

// This empty module with port declaration file causes synthesis tools to infer a black box for IP.
// The synthesis directives are for Synopsys Synplify support to prevent IO buffer insertion.
// Please paste the declaration into a Verilog source file or add the file as an additional source.
(* x_core_info = "blk_mem_gen_v8_4_2,Vivado 2018.3" *)
module data_ram(clka, wea, addra, dina, clkb, enb, addrb, doutb)
/* synthesis syn_black_box black_box_pad_pin="clka,wea[0:0],addra[5:0],dina[31:0],clkb,enb,addrb[5:0],doutb[31:0]" */;
  input clka;
  input [0:0]wea;
  input [5:0]addra;
  input [31:0]dina;
  input clkb;
  input enb;
  input [5:0]addrb;
  output [31:0]doutb;
endmodule
