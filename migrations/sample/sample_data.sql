DO $$
DECLARE
    i INT;
    j INT;
    k INT;
    current_user_id INT;
    workout_id INT;
    workout_type INT;
    workout_types TEXT[] := ARRAY['Running', 'Strength', 'HIIT', 'Yoga', 'Swimming', 'Cycling', 'CrossFit', 'Pilates'];
    
    -- Use separate arrays instead of multidimensional array
    running_exercises TEXT[] := ARRAY['Easy Run', 'Tempo Run', 'Interval Training', 'Hill Repeats', 'Long Run'];
    strength_exercises TEXT[] := ARRAY['Bench Press', 'Squats', 'Deadlift', 'Shoulder Press', 'Pull-ups', 'Rows', 'Leg Press', 'Bicep Curls'];
    hiit_exercises TEXT[] := ARRAY['Burpees', 'Mountain Climbers', 'Jump Squats', 'Push-ups', 'High Knees', 'Jumping Jacks'];
    yoga_exercises TEXT[] := ARRAY['Sun Salutation', 'Warrior Poses', 'Balance Poses', 'Seated Poses', 'Inversions'];
    swimming_exercises TEXT[] := ARRAY['Freestyle', 'Backstroke', 'Breaststroke', 'Butterfly', 'Drills'];
    cycling_exercises TEXT[] := ARRAY['Endurance Ride', 'Hill Climbs', 'Sprints', 'Recovery Ride', 'Intervals'];
    crossfit_exercises TEXT[] := ARRAY['Box Jumps', 'Kettlebell Swings', 'Wall Balls', 'Thrusters', 'Double Unders', 'Muscle-ups'];
    pilates_exercises TEXT[] := ARRAY['The Hundred', 'Roll Up', 'Leg Circles', 'Rolling Like a Ball', 'Single Leg Stretch'];
    
    locations TEXT[] := ARRAY['Home', 'Gym', 'Park', 'Track', 'Pool', 'Studio', 'Beach', 'Mountains'];
    adjectives TEXT[] := ARRAY['Intense', 'Easy', 'Challenging', 'Relaxing', 'Quick', 'Lengthy', 'Morning', 'Evening', 'Weekend', 'Lunchtime'];
    workout_title TEXT;
    workout_description TEXT;
    exercise_count INT;
    exercise_name TEXT;
    reps INT;
    weight DECIMAL;
    duration INT;
    exercise_array_size INT;
    chosen_exercise_array TEXT[];
BEGIN
    -- Loop to create 20 users
    FOR i IN 1..20 LOOP
        -- Insert user
        INSERT INTO users (username, email, password_hash, bio)
        VALUES (
            'user' || i, 
            'user' || i || '@example.com',
            '$2a$12$WpWpkbCXvx5s4LzA5o1Zj.tMTc0/VH1Jus9c3qV37h/oLIbQ12cRK', -- password '123'
            'Bio for user ' || i || '. Fitness enthusiast who enjoys ' || 
            workout_types[1 + (i % 8)] || ' and ' || workout_types[1 + ((i+3) % 8)] || '.'
        )
        RETURNING id INTO current_user_id;
        
        -- Create 20 workouts for each user
        FOR j IN 1..20 LOOP
            -- Choose workout type (0-7)
            workout_type := (i+j) % 8;
            
            -- Choose random workout type and adjective
            workout_title := adjectives[1 + (j % 10)] || ' ' || workout_types[1 + workout_type] || ' ' || locations[1 + ((i+j*2) % 8)];
            workout_description := 'A ' || adjectives[1 + ((j*3) % 10)] || ' ' || workout_types[1 + workout_type] || ' session at ' || locations[1 + ((i+j*2) % 8)];
            
            -- Insert workout
            INSERT INTO workouts (user_id, title, description, duration_minutes, calories_burned) 
            VALUES (
                current_user_id,
                workout_title,
                workout_description,
                -- Duration between 10 and 90 minutes
                10 + ((i+j) % 80),
                -- Calories between 50 and 800
                50 + ((i*j) % 750)
            )
            RETURNING id INTO workout_id;
            
            -- Determine how many exercise entries to add (2-6)
            exercise_count := 2 + (i + j) % 5;
            
            -- Select appropriate exercise array based on workout type
            CASE workout_type
                WHEN 0 THEN chosen_exercise_array := running_exercises;
                WHEN 1 THEN chosen_exercise_array := strength_exercises;
                WHEN 2 THEN chosen_exercise_array := hiit_exercises;
                WHEN 3 THEN chosen_exercise_array := yoga_exercises;
                WHEN 4 THEN chosen_exercise_array := swimming_exercises;
                WHEN 5 THEN chosen_exercise_array := cycling_exercises;
                WHEN 6 THEN chosen_exercise_array := crossfit_exercises;
                WHEN 7 THEN chosen_exercise_array := pilates_exercises;
            END CASE;
            
            exercise_array_size := array_length(chosen_exercise_array, 1);
            
            -- Add exercise entries
            FOR k IN 1..exercise_count LOOP
                -- Get appropriate exercise from the selected array
                exercise_name := chosen_exercise_array[1 + (k % exercise_array_size)];
                
                INSERT INTO workout_entries (
                    workout_id,
                    exercise_name,
                    sets,
                    reps,
                    duration_seconds,
                    weight,
                    notes,
                    order_index
                )
                VALUES (
                    workout_id,
                    exercise_name,
                    -- Sets: 1-5
                    1 + (k % 5),
                    -- Reps (only for strength/HIIT/CrossFit)
                    CASE WHEN workout_type IN (1, 2, 6) THEN 5 + (i*k % 15) ELSE NULL END,
                    -- Duration (for cardio/yoga)
                    CASE WHEN workout_type IN (0, 3, 4, 5, 7) THEN (30 + (j*k % 300)) ELSE NULL END,
                    -- Weight (only for strength/CrossFit)
                    CASE WHEN workout_type IN (1, 6) THEN 10.0 + (i*k % 90)::DECIMAL ELSE NULL END,
                    -- Notes
                    'Notes for exercise ' || k || ' of workout ' || j || ' by user' || i,
                    -- Order
                    k
                );
            END LOOP;
        END LOOP;
    END LOOP;
END $$;
