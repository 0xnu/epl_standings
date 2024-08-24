-- EPL Data Database
CREATE DATABASE `epl_data` DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE football_table (
    uuid            VARCHAR(36) PRIMARY KEY,
    position        VARCHAR(10),
    team_name       VARCHAR(100),
    played          VARCHAR(10),
    won             VARCHAR(10),
    drawn           VARCHAR(10),
    lost            VARCHAR(10),
    goals_for       VARCHAR(10),
    goals_against   VARCHAR(10),
    goal_difference VARCHAR(10),
    points          VARCHAR(10),
    form            VARCHAR(50)
);
