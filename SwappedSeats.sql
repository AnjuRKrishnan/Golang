Create table seat(
   id int,
   student varchar(225)
   );
   
   INSERT INTO seat (id, student)
VALUES 
    (1, 'Abbot'),
    (2, 'Doris'),
    (3, 'Emerson'),
    (4, 'Green'),
    (5, 'Jeames');
    
    select * from seat;
    
    WITH SwappedSeats AS (
    SELECT 
        CASE
            WHEN id % 2 = 1 AND id = (SELECT MAX(id) FROM Seat) THEN id
            WHEN id % 2 = 1 THEN id + 1
            ELSE id - 1
        END AS new_id,
        student
    FROM Seat
)
SELECT * FROM SwappedSeats order by new_id;
