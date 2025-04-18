with Ada.Text_IO; use Ada.Text_IO;
with Ada.Numerics.Float_Random; use Ada.Numerics.Float_Random;
with Random_Seeds; use Random_Seeds;
with Ada.Real_Time; use Ada.Real_Time;

-- proc running the simulation
procedure  Travelers is

-- -- Configuration:

-- moving on the board
   Nr_Of_Travelers : constant Integer := 15;

   Min_Steps : constant Integer := 10;
   Max_Steps : constant Integer := 100;

   Min_Delay : constant Duration := 0.01;
   Max_Delay : constant Duration := 0.05;

-- 2D Board Dimensions (with torus topology)
   Board_Width  : constant Integer := 15;
   Board_Height : constant Integer := 15;

-- Timing
   Start_Time : Time := Clock;  -- global starting time

-- Random seeds for the tasks' random number generators
   Seeds : Seed_Array_Type(1..Nr_Of_Travelers) := Make_Seeds(Nr_Of_Travelers);

-- -- Types, procedures and functions:

   type Position_Type is record
      X: Integer range 0 .. Board_Width - 1;
      Y: Integer range 0 .. Board_Height - 1;
   end record;

   procedure Move_Down(Position: in out Position_Type) is
   begin
      Position.Y := (Position.Y + 1) mod Board_Height;
   end Move_Down;

   procedure Move_Up(Position: in out Position_Type) is
   begin
      Position.Y := (Position.Y + Board_Height - 1) mod Board_Height;
   end Move_Up;

   procedure Move_Right(Position: in out Position_Type) is
   begin
      Position.X := (Position.X + 1) mod Board_Width;
   end Move_Right;

   procedure Move_Left(Position: in out Position_Type) is
   begin
      Position.X := (Position.X + Board_Width - 1) mod Board_Width;
   end Move_Left;

   type Trace_Type is record
      Time_Stamp:  Duration;
      Id : Integer;
      Position: Position_Type;
      Symbol: Character;
   end record;

   type Trace_Array_type is array(0 .. Max_Steps) of Trace_Type;

   type Traces_Sequence_Type is record
      Last: Integer := -1;
      Trace_Array: Trace_Array_type ;
   end record;

   procedure Print_Trace(Trace : Trace_Type) is
      Symbol : String := (' ', Trace.Symbol);
   begin
      Put_Line(
        Duration'Image(Trace.Time_Stamp) & " " &
        Integer'Image(Trace.Id) & " " &
        Integer'Image(Trace.Position.X) & " " &
        Integer'Image(Trace.Position.Y) & " " &
        (' ', Trace.Symbol)
);
   end Print_Trace;

   procedure Print_Traces(Traces : Traces_Sequence_Type) is
   begin
      for I in 0 .. Traces.Last loop
        Print_Trace(Traces.Trace_Array(I));
      end loop;
   end Print_Traces;

   task Printer is
      entry Report(Traces : Traces_Sequence_Type);
   end Printer;

   task body Printer is
   begin
      for I in 1 .. Nr_Of_Travelers loop -- range for TESTS !!!
         accept Report(Traces : Traces_Sequence_Type) do
           Print_Traces(Traces);
         end Report;
      end loop;
   end Printer;

   type Traveler_Type is record
      Id: Integer;
      Symbol: Character;
      Position: Position_Type;
   end record;

   task type Traveler_Task_Type is
      entry Init(Id: Integer; Seed: Integer; Symbol: Character);
      entry Start;
   end Traveler_Task_Type;

   task body Traveler_Task_Type is
      G : Generator;
      Traveler : Traveler_Type;
      Time_Stamp : Duration;
      Nr_of_Steps: Integer;
      Traces: Traces_Sequence_Type;

      procedure Store_Trace is
      begin
         Traces.Last := Traces.Last + 1;
         Traces.Trace_Array(Traces.Last) := (
            Time_Stamp => Time_Stamp,
            Id => Traveler.Id,
            Position => Traveler.Position,
            Symbol => Traveler.Symbol
         );
      end Store_Trace;

      procedure Make_Step is
         N : Integer;
      begin
         N := Integer(Float'Floor(4.0 * Random(G)));
         case N is
            when 0 =>
            Move_Up(Traveler.Position);
            when 1 =>
            Move_Down(Traveler.Position);
            when 2 =>
            Move_Left(Traveler.Position);
            when 3 =>
            Move_Right(Traveler.Position);
            when others =>
            Put_Line(" ?????????????? " & Integer'Image(N));
         end case;
      end Make_Step;

      function Get_Random_Delay return Duration is
      begin
         return Min_Delay + (Max_Delay-Min_Delay) * Duration(Random(G));
      end Get_Random_Delay;

   begin

      accept Init(Id: Integer; Seed: Integer; Symbol: Character) do
         Reset(G, Seed);
         Traveler.Id := Id;
         Traveler.Symbol := Symbol;

         -- Random initial position:
         Traveler.Position := (
            X => Integer(Float'Floor(Float(Board_Width)  * Random(G))),
            Y => Integer(Float'Floor(Float(Board_Height) * Random(G)))
         );

         Store_Trace; -- store starting position

         -- Number of steps to be made by the traveler
         Nr_of_Steps := Min_Steps + Integer(Float(Max_Steps - Min_Steps) * Random(G));
         -- Time_Stamp of initialization
         Time_Stamp := To_Duration (Clock - Start_Time); -- reads global clock
      end Init;

      -- wait for remaining tasks to commence:
      accept Start do
         null;
      end Start;

      for Step in 0 .. Nr_of_Steps loop
         delay Get_Random_Delay;
         -- do action ...
         Make_Step;
         Store_Trace;
         Time_Stamp := To_Duration (Clock - Start_Time); -- reads global clock
      end loop;

      -- schedule a report
      Printer.Report(Traces);

   end Traveler_Task_Type;


-- variables declarations for main task

   -- number of travelers
   Travel_Tasks: array (0 .. Nr_Of_Travelers-1) of Traveler_Task_Type;
   Symbol : Character := 'A';
begin

   Put_Line(
      "timestamp "&
      "| no. travelers" &" "&
      "| width" &" "&
      "| height" &" "&
      "| symbol"
   );
   Put_Line(
      "... |"&
      Integer'Image(Nr_Of_Travelers) &" |"&
      Integer'Image(Board_Width) &" |"&
      Integer'Image(Board_Height) &" |"&
      "..."
   );


   for I in Travel_Tasks'Range loop
      Travel_Tasks(I).Init(I, Seeds(I+1), Symbol);
      Symbol := Character'Succ(Symbol);
   end loop;

   for I in Travel_Tasks'Range loop
      Travel_Tasks(I).Start;
   end loop;

end Travelers;

