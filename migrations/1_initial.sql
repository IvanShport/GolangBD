-- +migrate Up

CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE user_profile (
    nickname citext PRIMARY KEY,
    email citext UNIQUE NOT NULL,
    fullname text,
    about text
);

CREATE TABLE forum (
    posts integer DEFAULT 0,
    forum_slug citext PRIMARY KEY,
    threads integer DEFAULT 0,
    forum_title text NOT NULL,
    forum_user citext REFERENCES user_profile NOT NULL
);

CREATE TABLE thread (
    thread_author citext REFERENCES user_profile NOT NULL,
    thread_created timestamp with time zone DEFAULT now(),
    thread_forum citext references forum NOT NULL,
    thread_id serial PRIMARY KEY,
    thread_message text NOT NULL,
    thread_slug citext UNIQUE,
    thread_title text NOT NULL,
    votes integer DEFAULT 0
);

CREATE TABLE post (
    post_author citext REFERENCES user_profile NOT NULL,
    post_created timestamp with time zone DEFAULT now(),
    post_forum citext references forum NOT NULL,
    post_id serial PRIMARY KEY,
    isEdited boolean DEFAULT FALSE NOT NULL,
    post_message text NOT NULL,
    parent integer default 0,
    post_thread integer REFERENCES thread NOT NULL,
    path integer ARRAY,
    founder integer DEFAULT 0
);

CREATE TABLE vote (
    nickname citext REFERENCES user_profile NOT NULL,
    thread integer REFERENCES thread NOT NULL,
    voice integer NOT NULL CHECK (voice IN (-1, 1)),
    CONSTRAINT uniq_vote UNIQUE (nickname, thread)
);

CREATE TABLE forum_users (
    nickname citext REFERENCES user_profile,
    forum citext REFERENCES forum,
    CONSTRAINT uniq_forum_users UNIQUE (nickname, forum)
);


-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION add_thread() RETURNS TRIGGER AS $add_thread$
    BEGIN
        UPDATE forum SET threads = threads + 1 WHERE forum_slug = NEW.thread_forum;
        RETURN NEW;
    END;
$add_thread$ LANGUAGE plpgsql;

CREATE TRIGGER add_thread AFTER INSERT ON thread
FOR EACH ROW EXECUTE PROCEDURE add_thread();
-- +migrate StatementEnd


-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION add_vote() RETURNS TRIGGER AS $add_vote$
    BEGIN
        IF (TG_OP = 'INSERT') THEN
            UPDATE thread SET votes = votes + NEW.voice WHERE thread_id = NEW.thread;
            RETURN NEW;
        ELSIF (TG_OP = 'UPDATE') THEN
            IF OLD.voice <> NEW.voice THEN
                UPDATE thread SET votes = votes + NEW.voice * 2 WHERE thread_id = NEW.thread;
            END IF;
            RETURN NEW;
        END IF;
        RETURN NULL;
    END;
$add_vote$ LANGUAGE plpgsql;
-- +migrate StatementEnd

CREATE TRIGGER add_vote AFTER INSERT OR UPDATE ON vote
FOR EACH ROW EXECUTE PROCEDURE add_vote();