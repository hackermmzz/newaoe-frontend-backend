module load_use_hazard(
    input  [4:0] id_rs,
    input  [4:0] id_rt,
    input  [4:0] id_ex_rt,
    input        id_ex_memread,
    output       stall,
    output       bubble
);
    wire rs_conflict = (id_ex_rt == id_rs);
    wire rt_conflict = (id_ex_rt == id_rt);
    wire load_dependency = id_ex_memread && (rs_conflict || rt_conflict);
    assign stall  = load_dependency;
    assign bubble = load_dependency;
endmodule
