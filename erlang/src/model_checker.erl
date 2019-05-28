%%%-------------------------------------------------------------------
%%% @author tyryhryrhyhr
%%% @copyright (C) 2019, <COMPANY>
%%% @doc
%%%
%%% @end
%%% Created : 08. May 2019 8:26 AM
%%%-------------------------------------------------------------------
-module(model_checker).
-author("tyryhryrhyhr").

%% API

-export([measureMC/0, startMC/0, explore/3, create_location/5, create_location/1, create_edge/2, create_edge/5,
        create_edge/6]).

%setup model, create processes and start exploring

measureMC()->

eprof:start(),
  eprof:start_profiling([self()]),
statistics(runtime),
statistics(wall_clock),

startMC(),

{_, Time1} = statistics(runtime),
{_, Time2} = statistics(wall_clock),
U1 = Time1 * 1000,
U2 = Time2 * 1000,
io:format("Code time=~p (~p) microseconds~n",
[U1,U2]),
  eprof:stop_profiling(),
  eprof:analyze(total).

startMC() ->

  io:fwrite("Startingmof\n"),

  Model = [create_process(derpa),create_process(derpb),create_process(derpc)],
  {Starting_vector,Processlist} = lists:unzip(Model),
  io:fwrite("Starting Vector : ~w vectors~w\n",[Starting_vector,digraph:edges(hd(Processlist))]),
  Waiting_queue = queue:in([Starting_vector,[[0],[0],[0],[0]]],queue:new()),
  io:fwrite("Waiting queue: ~w\n",[Waiting_queue]),
  explore(Waiting_queue, queue:new(), Processlist).

create_process(second) ->
  Graph = digraph:new(),
  %Vertices
  V1=digraph:add_vertex(Graph, "A"),
  V2=digraph:add_vertex(Graph, "B"),
  V3=digraph:add_vertex(Graph, "C"),
  V4=digraph:add_vertex(Graph, "D"),
  V5=digraph:add_vertex(Graph, "E"),
  V6=digraph:add_vertex(Graph, "F"),
  V7=digraph:add_vertex(Graph, "G"),
  V8=digraph:add_vertex(Graph, "H"),
  V9=digraph:add_vertex(Graph, "I"),
  %Edges
  E1=digraph:add_edge(Graph,V1,V2),
  E2=digraph:add_edge(Graph,V2,V3),
  E3=digraph:add_edge(Graph,V2,V5),
  E4=digraph:add_edge(Graph,V3,V1),
  E5=digraph:add_edge(Graph,V4,V6),
  E6=digraph:add_edge(Graph,V5,V6),
  E7=digraph:add_edge(Graph,V5,V7),
  E8=digraph:add_edge(Graph,V3,V8),
  E9=digraph:add_edge(Graph,V7,V9),
  E10=digraph:add_edge(Graph,V8,V9),
  io:fwrite("~w\n",[Graph]),
  Process={V1,Graph},
  %io:fwrite("Vertices: ~w\n",[digraph:vertices(Graph)]),
  %io:fwrite("Edges: ~w\n",[digraph:edges(Graph)]),
  %io:fwrite("Outgoing from start: ~w\n",[digraph:out_edges(Graph,V1)]),
  Process;

create_process(first) ->
  Graph = digraph:new(),
  %Vertices
  V1=digraph:add_vertex(Graph, "A"),
  V2=digraph:add_vertex(Graph, "B"),
  V3=digraph:add_vertex(Graph, "C"),
  %Edges
  E1=digraph:add_edge(Graph,V1,V2,[2,plus,1]),
  E2=digraph:add_edge(Graph,V2,V3,[send,1]),
  E3=digraph:add_edge(Graph,V1,V3),
  %io:fwrite("~w\n",[Graph]),
  Process={V1,Graph},
  %io:fwrite("Vertices: ~w\n",[digraph:vertices(Graph)]),
  %io:fwrite("Edges: ~w\n",[digraph:edges(Graph)]),
  %io:fwrite("Outgoing from start: ~w\n",[digraph:out_edges(Graph,V1)]),
  Process;

create_process(derpa) ->
  Graph = digraph:new(),
  %Vertices
  V1=digraph:add_vertex(Graph, "A"),
  %Edges
  E1=digraph:add_edge(Graph,V1,V1,[1,plus,1]),
  %io:fwrite("~w\n",[Graph]),
  Process={V1,Graph},
  %io:fwrite("Vertices: ~w\n",[digraph:vertices(Graph)]),
  %io:fwrite("Edges: ~w\n",[digraph:edges(Graph)]),
  %io:fwrite("Outgoing from start: ~w\n",[digraph:out_edges(Graph,V1)]),
  Process;

create_process(derpb) ->
  Graph = digraph:new(),
  %Vertices
  V1=digraph:add_vertex(Graph, "A"),
  %Edges
  E1=digraph:add_edge(Graph,V1,V1,[2,plus,1]),
  io:fwrite("~w\n",[Graph]),
  Process={V1,Graph},
  %io:fwrite("Vertices: ~w\n",[digraph:vertices(Graph)]),
  %io:fwrite("Edges: ~w\n",[digraph:edges(Graph)]),
  %io:fwrite("Outgoing from start: ~w\n",[digraph:out_edges(Graph,V1)]),
  Process;

create_process(derpc) ->
  Graph = digraph:new(),
  %Vertices
  V1=digraph:add_vertex(Graph, "A"),
  %Edges
  E1=digraph:add_edge(Graph,V1,V1,[3,plus,1]),
  io:fwrite("~w\n",[Graph]),
  Process={V1,Graph},
  %io:fwrite("Vertices: ~w\n",[digraph:vertices(Graph)]),
  %io:fwrite("Edges: ~w\n",[digraph:edges(Graph)]),
  %io:fwrite("Outgoing from start: ~w\n",[digraph:out_edges(Graph,V1)]),
  Process;


create_process(ohbabydouble) ->
  Graph = digraph:new(),
  %Vertices
  V1=digraph:add_vertex(Graph, "A"),
  %Edges
  E1=digraph:add_edge(Graph,V1,V1,[1,plus,1]),
  E2=digraph:add_edge(Graph,V1,V1,[2,plus,1]),
  io:fwrite("~w\n",[Graph]),
  Process={V1,Graph},
  %io:fwrite("Vertices: ~w\n",[digraph:vertices(Graph)]),
  %io:fwrite("Edges: ~w\n",[digraph:edges(Graph)]),
  %io:fwrite("Outgoing from start: ~w\n",[digraph:out_edges(Graph,V1)]),
  Process;

create_process(empty) ->
  Graph = digraph:new(),
  %Vertices
  V1=digraph:add_vertex(Graph, "A"),
  %Edges
  io:fwrite("~w\n",[Graph]),
  Process={V1,Graph},
  %io:fwrite("Vertices: ~w\n",[digraph:vertices(Graph)]),
  %io:fwrite("Edges: ~w\n",[digraph:edges(Graph)]),
  %io:fwrite("Outgoing from start: ~w\n",[digraph:out_edges(Graph,V1)]),
  Process;

create_process(synca) ->
  Graph = digraph:new(),
  %Vertices
  V1=digraph:add_vertex(Graph, "A"),
  V2=digraph:add_vertex(Graph, "B"),
  %Edges
  E1=digraph:add_edge(Graph,V1,V2,[send,1]),
  io:fwrite("~w\n",[Graph]),
  Process={V1,Graph},
  %io:fwrite("Vertices: ~w\n",[digraph:vertices(Graph)]),
  %io:fwrite("Edges: ~w\n",[digraph:edges(Graph)]),
  %io:fwrite("Outgoing from start: ~w\n",[digraph:out_edges(Graph,V1)]),
  Process;

create_process(syncb) ->
  Graph = digraph:new(),
  %Vertices
  V1=digraph:add_vertex(Graph, "A"),
  V2=digraph:add_vertex(Graph, "B"),
  %Edges
  E1=digraph:add_edge(Graph,V1,V2,[recieve,1]),
  io:fwrite("~w\n",[Graph]),
  Process={V1,Graph},
  %io:fwrite("Vertices: ~w\n",[digraph:vertices(Graph)]),
  %io:fwrite("Edges: ~w\n",[digraph:edges(Graph)]),
  %io:fwrite("Outgoing from start: ~w\n",[digraph:out_edges(Graph,V1)]),
  Process.



%create node from name and invariants if any

create_location(Name,invariant, Var, Operator, Value) -> 0.

create_location(Name) -> 0.

%create edge from starting node, target node, guards and updates

create_edge(Source, Target) -> 0.


create_edge(Source, Target, guard, Var, Operator, Value) -> 0;


create_edge(Source, Target, update, Var, Operator, Value) -> 0.

create_edge(Source, Target, sync, Channel, Value) -> 0.

%Given a waiting queue and an explored queue recursively explores the waiting queue until it is empty.

explore({[],[]}, Explored_queue, _) ->
  io:fwrite("we did it patrick"),
  io:fwrite("Waiting queue is empty, total explored: ~w\n",[queue:len(Explored_queue)]),
  true;

explore(Waiting_queue, Explored_queue, Processlist) ->

  %Exploration loop, continously checks top queue element from waiting queue until it is empty, and adds vectors from explorig unless they are duplica<<<<<<<<<<<<<<<tes.
  %io:fwrite("Explore Initial queue: ~w\n",[Waiting_queue]),
  {{value,[Location_vector,Variable_vector]},Waiting_mid} = queue:out(Waiting_queue),
  %io:fwrite("exploring ~w\n",[Waiting_mid]),
  Newvectors=accumulate_vertices(Location_vector,Variable_vector,Processlist,Processlist,[],[]),
  %io:fwrite("Explored ~w and found~w \n",[[Location_vector,Variable_vector],Newvectors]),
  Waiting_new = enqueue_new_vectors(Waiting_mid, Explored_queue, Newvectors),
  %io:fwrite("Exploration complete ~w\n",[Explored_queue]),
  Explored_new=queue:in([Location_vector,Variable_vector],Explored_queue),
  %io:fwrite("Exploration complete, new Explored queue: ~w\n",[Explored_new]),
  explore(Waiting_new,Explored_new,Processlist).
  %insert exploration step here
  %Create threading and supervision here as well
  %io:fwrite("kkk").
  %A model is a list of processes.
  %A Process is a tuple consisting of a starting vector, a labelled digraph and a list of local variables. Digraph labels on vertex are records of "A",operator,"B" defining invariants.
  % Digraph labels on edges are also record of identical format but can be guards, updates or syncs

create_machine(Cores) ->
  [].

create_thread(Waiting_queue,Explored_queue,Processlist,Adresslist)->
  [].

%thread_readloop({[],[]},Explored_queue,Processlist,Addresslist)->
  %receive
  %  Veclist -> []
  %end,
  %Waiting_recieved = enqueue_new_vectors(Waiting_queue,Explored_queue,Veclist),
  %thread_readloop(Waiting_recieved,Explored_queue,Processlist,Addresslist);


%thread_readloop(Waiting_queue,Explored_queue,Processlist,Addresslist)->
  %{_,Veclist} = recieve,
  %Waiting_recieved = enqueue_new_vectors(Waiting_queue,Explored_queue,Veclist),
  %{{value,[Location_vector,Variable_vector]},Waiting_new} = queue:out(Waiting_recieved),
  %Newvectors=accumulate_vertices(Location_vector,Variable_vector,Processlist,Processlist,[],[]),
  %hash_distribution(Newvectors,Addresslist),
  %Explored_new=queue:in([Location_vector,Variable_vector],Explored_queue),
  %thread_readloop(Waiting_new,Explored_new,Processlist,Addresslist).




enqueue_new_vectors(Queue,_,[]) ->
  %io:fwrite("wef ~w\n",[Queue]),
  Queue;
enqueue_new_vectors(Waiting, Explored,[H|T]) ->
  %io:fwrite("Adding ~w to queue~w\n",[H,Waiting]),
  Vectoralreadyinqueue = queue:member(H,Waiting),

  case Vectoralreadyinqueue of

    false ->
      Vectoralreadyexplored = queue:member(H,Explored),
      %Vectoralreadyexplored = false,
      case Vectoralreadyexplored of
        false  ->
          Q2=queue:in(H,Waiting),
          %io:fwrite("Adding ~w to queue~w\n",[H,Q2]),
          enqueue_new_vectors(Q2,Explored,T);
      true ->
        %io:fwrite("vector already explored"),
        enqueue_new_vectors(Waiting,Explored,T)
end;
      %io:fwrite("aaa"),

    true ->
      %io:fwrite("vector already in queue"),
      enqueue_new_vectors(Waiting,Explored,T)

  end.


accumulate_vertices([],_,[],_,_,Acc)->
  %Case checks if empty
  Acc;
accumulate_vertices(Location_vector,Variable_vector,[Processhead|Processtail],Model,Acchead,Acc)->
  %io:fwrite("edgeget ~w ~w \n",[Location_vector,digraph:out_edges(Processhead,hd(Location_vector))]),
  Newvec=accumulate_edges(Location_vector,Variable_vector,digraph:out_edges(Processhead,hd(Location_vector)),Processhead,Processtail,[]),
  %io:fwrite("Newvec~w\n",[Newvec]),
  Newacc=list_assembler(Acchead,Newvec)++Acc,
  %io:fwrite("explore steep accumulator: ~w\n",[Newacc]),
  %io:fwrite("explore steep accumulator: ~w\n",[Acchead]),
  accumulate_vertices(tl(Location_vector),Variable_vector,Processtail,Model,[hd(Location_vector)|Acchead],Newacc).


list_assembler(_,[])->
  [];
list_assembler([],List)->
  %io:fwrite("Assembled ~w\n",[List]),
  List;
list_assembler([Prefix|Rec],List)->
  %io:fwrite("Assembling ~w~w\n",[List,Prefix]),
  Newlist=prefix_rec(Prefix,List,[]),
  list_assembler(Rec,Newlist).

prefix_rec(_,[],Acc)->
  Acc;
prefix_rec(Var,List,Acc)->
  Newacc=[[[Var|hd(hd(List))]|tl(hd(List))]|Acc],
  %io:fwrite("Var~w, List~w, Acc~w\n",[Var,List,Acc]),
  prefix_rec(Var,tl(List),Newacc).




accumulate_edges(_,_,[],_,_,Acc) ->
  Acc;
accumulate_edges(Location_vector,Variable_vector,Edgelist,Process,Processtail,Acc) ->
  %Finds all state vectors across multiple edges
  %io:fwrite("Taking edge: ~w~w\n",[Location_vector,hd(Edgelist)]),
  Newacc=take_edge(Location_vector,Variable_vector,hd(Edgelist),Process,Processtail),
  %io:fwrite("New vector found: ~w\n",[Newacc]),
  case Newacc of
    [] -> %io:fwrite("New vector was empty: ~w\n",[Newacc]),
      accumulate_edges(Location_vector,Variable_vector,tl(Edgelist),Process,Processtail,Acc);
    _ -> %io:fwrite("adding new vector: ~w\n",[Newacc]),
      accumulate_edges(Location_vector,Variable_vector,tl(Edgelist),Process,Processtail,[Newacc|Acc])


  end.
  %accumulate_edges(Location_vector,Variable_vector,tl(Edgelist),Process,Model,Newacc).

%accumulate_processes(Location_vector,Variable_vector,[],Model,Acc) ->
%  Acc;
%
%accumulate_processes(Location_vector,Variable_vector,Process_list,Model,Acc) ->
%  Newacc = [explore_vertex(Location_vector,Variable_vector,hd(Process_list),Model)|list_assembler(hd(Location_vector),Acc)],
%  io:fwrite("accumulate processes found vectors ~w\n",[Newacc]),
%  accumulate_processes(tl(Location_vector),Variable_vector,tl(Process_list),Model,Newacc).


take_edge(Location_vector,Variable_vector,Edge,Process,Processtail)->
  {Edge,_,Dest,Labellist}=digraph:edge(Process,Edge),
  %io:fwrite("Taking edge ~w with labels~w and location vector~w\n",[Edge,Labellist,Location_vector]),
  case Labellist of
    [] ->
      %io:fwrite("No Label ~w\n",[[Dest|tl(Location_vector)],Variable_vector]),
      [[Dest|tl(Location_vector)],Variable_vector];
    [send,Channel] ->
      Syncvectors=syncrec(send,tl(Location_vector),Variable_vector,Processtail,Channel),
      list_assembler([Dest],Syncvectors);
    [recieve,Channel] ->
      Syncvectors=syncrec(recieve,tl(Location_vector),Variable_vector,Processtail,Channel),
      list_assembler([Dest],Syncvectors);
    [A,set,B] ->
      %io:fwrite("Assign label"),
      Newvarvec=lists:flatten([lists:sublist(Variable_vector,A-1)++[B]|lists:nthtail(A,Variable_vector)]),
      [[Dest|tl(Location_vector)],Newvarvec];
    [A,plus,B] ->
      %io:fwrite("Addition label ~w\n",[hd(lists:nth(A,Variable_vector))]),
      N=hd(lists:nth(A,Variable_vector))+B,
      %io:fwrite("Newvalue~w\n",[N]),
      if
        N =<(50) ->
          %io:fwrite("mof"),
          Newvarvec=lists:sublist(Variable_vector,A-1)++[[N]]++lists:nthtail(A,Variable_vector),
          [[Dest|tl(Location_vector)],Newvarvec];

        true ->
          %io:fwrite("nono"),
          []
      end;

    [A,minus,B] ->
      Newvarvec=[lists:sublist(Variable_vector,A-1)++[hd(lists:nth(A,Variable_vector))-B]|lists:nthtail(A,Variable_vector)],
      [[Dest|tl(Location_vector)],Newvarvec];
    [A,lessthan,B] ->
      Guard=hd(lists:nth(A,Variable_vector))<B,
      if Guard ->
        %io:fwrite("Guard Succeed\n"),
        Result=[[Dest|tl(Location_vector)],Variable_vector];
        true ->
          %io:fwrite("Guard Failed\n"),
          Result=[]
      end,
      Result;
    [A,greaterthan,B] ->
      Guard=hd(lists:nth(A,Variable_vector))<B,
      if Guard ->
        %io:fwrite("Guard Succeed\n"),
        Result=[[Dest|tl(Location_vector)],Variable_vector];
        true ->
          %io:fwrite("Guard Failed\n"),
          Result=[]
      end,
      Result;
    _else -> io:fwrite("Got a label i cant read")
  end.

syncrec(recieve,Location_vector,Variable_vector,Processlist,Channel) ->
  %io:fwrite("\n SYNC TEST RECIEVE ~w ~w ~w\n\n",[Location_vector,Variable_vector,Processlist]),
  syncrec_process(send,Location_vector,Variable_vector,Processlist,Channel,[],[]);

syncrec(send,Location_vector,Variable_vector,Processlist,Channel) ->
  %io:fwrite("\n SYNC TEST SEND ~w ~w ~w\n\n",[Location_vector,Variable_vector,Processlist]),
  syncrec_process(recieve,Location_vector,Variable_vector,Processlist,Channel,[],[]).

syncrec_edge(_,_,_,[],_,_,Acc)->
  %io:fwrite("Sync vector results: ~w\n",[Acc]),
  Acc;
syncrec_edge(Location_tail,Variable_vector,Sendrec,Edgelist,Process,Channel,Acc)->
  {_,_,V2,Label}=digraph:edge(Process,hd(Edgelist)),
  %io:fwrite("synccheck label~w match~w destination~w\n",[Label,Sendrec,V2]),
  case {Label,Acc} of
    {[Sendrec,Channel],[]} ->
      %io:fwrite("v2~w location tail~w variable vector~w acc ~w\n",[V2,Location_tail,Variable_vector,Acc]),
      %nth element location vector,variable vector
      Newacc=[[V2|Location_tail],Variable_vector],
      %io:fwrite("Newacc ~w cons ~w\n",[Newacc,Newacc|Acc]),
      syncrec_edge(Location_tail,Variable_vector,Sendrec,tl(Edgelist),Process,Channel,Newacc);
    {[Sendrec,Channel],Acc} ->
      %io:fwrite("v2~w location tail~w variable vector~w acc ~w\n",[V2,Location_tail,Variable_vector,Acc]),
      %nth element location vector,variable vector
      Newacc=[[V2|Location_tail],Variable_vector],
      %io:fwrite("Newacc ~w cons ~w\n",[Newacc,Newacc|Acc]),
      syncrec_edge(Location_tail,Variable_vector,Sendrec,tl(Edgelist),Process,Channel,[Newacc|Acc]);
    _ ->
      syncrec_edge(Location_tail,Variable_vector,Sendrec,tl(Edgelist),Process,Channel,Acc)
  end.

syncrec_process(_,[],_,_,_,_,Acc) ->
  Acc;

syncrec_process(Sendrec,Location_vector,Variable_vector,Processlist,Channel,Acchead,Acc) ->
  Edgelist=digraph:out_edges(hd(Processlist),hd(Location_vector)),
  %Store all subvectors that can synchronize at current vertex in syncvectors
  Syncvectors=syncrec_edge(tl(Location_vector),Variable_vector,Sendrec,Edgelist,hd(Processlist),Channel,[]),
  %io:fwrite("syncvectors~w acchead~w\n",[Syncvectors,Acchead]),
  Newacc=list_assembler(Acchead,Syncvectors)++Acc,
  %io:fwrite("Newacc~w Acchead~w Acc~w\n",[Newacc,Acchead,Acc]),
  syncrec_process(Sendrec,tl(Location_vector),Variable_vector,tl(Processlist),Channel,[hd(Location_vector)|Acchead],Newacc).