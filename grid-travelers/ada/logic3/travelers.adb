with Ada.Text_IO; use Ada.Text_IO;
with Ada.Numerics.Float_Random; use Ada.Numerics.Float_Random;
with Random_Seeds; use Random_Seeds;
with Ada.Real_Time; use Ada.Real_Time;
with Ada.Characters.Handling; use Ada.Characters.Handling;


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

-- Exceptions
   Step_Error : exception;
   Semaphore_Overflow : exception;

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
      Direction: Integer range 0 .. 3;
   end record;

   protected type Semaphore_Type(size: Natural) is
      entry Take_Semaphore;
      entry Release_Semaphore;
   private
      capacity: Natural := size;
      load: Natural     := size;
   end Semaphore_Type;

   protected body Semaphore_Type is
      entry Take_Semaphore when load > 0 is
      begin
         load := load - 1;
      end Take_Semaphore;
      
      entry Release_Semaphore when load < capacity is
      begin
         load := load + 1;
      end Release_Semaphore;
   end Semaphore_Type;

   -- Synchronize travelers
   type Semaphore_Pool_Type is array (0 .. Board_Width - 1, 0 .. Board_Height - 1) of Semaphore_Type(1);
   Semaphore_Pool: Semaphore_Pool_Type;

   task type Traveler_Task_Type is
      entry Init(Id: Integer; Seed: Integer; Symbol: Character; Position: Position_Type);
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

      procedure Make_Step(step_timeout: in Duration) is
         newPos : Position_Type := Traveler.Position; -- deep copy
         success : Boolean;
      begin
         case Traveler.Direction is
            when 0 =>
            Move_Up(newPos);
            when 1 =>
            Move_Down(newPos);
            when 2 =>
            Move_Left(newPos);
            when 3 =>
            Move_Right(newPos);
            when others =>
            Put_Line("Error in Move procedure for traveler " &
                     Integer'Image(Traveler.Id));
         end case;

         -- check if new position is valid
         select
            Semaphore_Pool(newPos.X, newPos.Y).Take_Semaphore;
         or
            delay step_timeout;
            raise Step_Error; 
         end select;

         -- if valid: make the move
         Semaphore_Pool(Traveler.Position.X, Traveler.Position.Y).Release_Semaphore;
         Traveler.Position := newPos;

      end Make_Step;

      function Get_Random_Delay return Duration is
      begin
         return Min_Delay + (Max_Delay-Min_Delay) * Duration(Random(G));
      end Get_Random_Delay;

      random_delay : Duration;
   begin

      accept Init(Id: in Integer; Seed: in Integer; Symbol: in Character; Position: in Position_Type) do
         Reset(G, Seed);
         Traveler.Id := Id;
         Traveler.Symbol := Symbol;

         Traveler.Position := Position;
         Traveler.Direction := Integer(Random(G) + 0.5); -- randomize direction
         -- move up or down
         if Traveler.Id mod 2 = 0 then
            -- or move left or right
            Traveler.Direction := Traveler.Direction + 2;
         end if;
         Semaphore_Pool(Traveler.Position.X, Traveler.Position.Y).Take_Semaphore;

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
         random_delay := Get_Random_Delay;
         delay random_delay;
         begin
            -- try to do action ...
            Make_Step(step_timeout => Max_Delay - random_delay);
            Store_Trace;
            Time_Stamp := To_Duration (Clock - Start_Time); -- reads global clock
         exception
            when Step_Error =>
               --  Put_Line("Timeout reached in step " & Integer'Image(Step));
               --  Put_Line("Traveler " & Integer'Image(Traveler.Id) & " stopped.");
               --  Put_Line("Position: " & Integer'Image(Traveler.Position.X) & " " &
                        --  Integer'Image(Traveler.Position.Y));
               --  Put_Line("Symbol: " & Traveler.Symbol);
               Traveler.Symbol := To_Lower(Traveler.Symbol); -- mark the traveler
               Store_Trace;
               exit;
            when Semaphore_Overflow =>
               Put_Line("Traveler " & Integer'Image(Traveler.Id) &
                        " internal semaphore error ");
               exit;
            when others =>
               Put_Line("Traveler " & Integer'Image(Traveler.Id) &
                        " unknown error");
               exit;
         end;
      end loop;

      -- schedule a report
      Printer.Report(Traces);

   end Traveler_Task_Type;

   function Get_Initial_Position(Index: Natural) return Position_Type is
   begin
      return (X => Index, Y => Index);
   end Get_Initial_Position;


-- variables declarations for main task
   Travel_Tasks: array (0 .. Nr_Of_Travelers-1) of Traveler_Task_Type;
   Symbol : Character := 'A';
   Initial_Positions : array (0 .. Nr_Of_Travelers-1) of Position_Type;
   Used_Initial_Positions : array (0 .. Board_Width, 0 .. Board_Height) of Boolean := (others => (others => False));
begin

   -- deprecated:
   --  Put_Line(
   --     "timestamp "&
   --     "| no. travelers" &" "&
   --     "| width" &" "&
   --     "| height" &" "&
   --     "| symbol"
   --  );
   --  Put_Line(
   --     "... |"&
   --     Integer'Image(Nr_Of_Travelers) &" |"&
   --     Integer'Image(Board_Width) &" |"&
   --     Integer'Image(Board_Height) &" |"&
   --     "..."
   --  );

   Put_Line("-1 15 15 15");

   for I in Travel_Tasks'Range loop
      Travel_Tasks(I).Init(I, Seeds(I+1), Symbol, Get_Initial_Position(I));
      Symbol := Character'Succ(Symbol);
   end loop;

   for I in Travel_Tasks'Range loop
      Travel_Tasks(I).Start;
   end loop;

end Travelers;

