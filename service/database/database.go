package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	// User Operations
	CheckIfUserExistsByUsername(username string) (bool, error)
	CheckIfUserExistsByUserId(userId int) (bool, error)
	GetUserId(usename string) (int, error)
	CreateUser(username string) error
	GetUserById(userId int) (User, error)
	GetUserByUsername(username string) (User, error)
	GetAllUsers() ([]User, error)
	UpdateUsername(username string, userId int) error
	UpdateUserPhotoUrl(url string, userId int) error

	// Conversation Operations
	GetPrivateConversationsByUserId(userId int) ([]ConversationId, error)
	GetGroupConversationsByUserId(userId int) ([]ConversationId, error)
	CheckIfConversationExistsByConversationId(conversationId int) (bool, error)
	GetConversationByConversationId(conversationId int, userId int) (Conversation, error)
	CreatePrivateConversation(userId_1 int, userId_2 int) (int, error)

	// Message Operations
	CheckIfMessageExistsByMessageId(messageId int) (bool, error)
	CreateMessage(timestamp string, senderId int, conversationId int, status string, mexType string, content string) error
	ReplyToAMessage(timestamp string, senderId int, conversationId int, status string, mexType string, content string, repliedTo int) error
	DeleteMessage(messageId int) error
	ForwardMessage(messageId int, senderId int, conversationId int, timestamp string) (Message, error)
	UpdateMessageStatus(messageId int, status string) error
	AddReceiver(userId int, messageId int) error
	AddReader(userId int, messageId int) error

	// Comment Operations
	CheckIfCommentExistsByCommentId(commentId int) (bool, error)
	CheckIfCommentExists(senderId int, messageId int) (bool, error)
	CheckIfIsUserComment(senderId int, commentId int) (bool, error)
	UpdateComment(senderId int, messageId int, reaction string, timestamp string) error
	AddComment(timestamp string, senderId int, messageId int, reaction string) error
	DeleteComment(commentId int) error

	// Group Operations
	CheckIfGroupExistsByGroupId(groupId int) (bool, error)
	CheckIfUserIsPartecipant(userId int, groupId int) (bool, error)
	CreateGroupConversation(name string, description string, photoUrl string) (int, error)
	AddUserToGroup(userId int, groupId int) error
	LeaveGroup(userId int, groupId int) error
	UpdateGroupName(groupId int, name string) error
	UpdateGroupDescription(groupId int, description string) error
	UpdateGroupPhotoUrl(url string, groupId int) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='User';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {

		// Create User Table
		sqlStmt := `CREATE TABLE User (
			userId INTEGER PRIMARY KEY AUTOINCREMENT, 
			username TEXT NOT NULL UNIQUE, 
			photoUrl TEXT NOT NULL);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating the user table: %w", err)
		}

		// Create Conversation Table
		sqlStmt = `CREATE TABLE Conversation (conversationId INTEGER PRIMARY KEY AUTOINCREMENT);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating the conversation table: %w", err)
		}

		// Create PrivateConversation Table
		sqlStmt = `CREATE TABLE PrivateConversation (
			conversationId INTEGER PRIMARY KEY,
			userId_1 INTEGER NOT NULL,
			userId_2 INTEGER NOT NULL,
			CHECK (userId_1 < userId_2),
			UNIQUE (userId_1, userId_2),
			FOREIGN KEY (conversationId) REFERENCES Conversation(conversationId),
			FOREIGN KEY (userId_1) REFERENCES User(userId),
			FOREIGN KEY (userId_2) REFERENCES User(userId));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating the private conversation table: %w", err)
		}

		// Create GroupConversation Table
		sqlStmt = `CREATE TABLE GroupConversation (
			groupId INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			photoUrl TEXT NOT NULL,
			conversationId INTEGER NOT NULL,
			FOREIGN KEY (conversationId) REFERENCES Conversation(conversationId));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating the group conversation table: %w", err)
		}

		// Create Message Table
		sqlStmt = `CREATE TABLE Message (
			messageId INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp TEXT NOT NULL,
			senderId INTEGER NOT NULL, 
			conversationId INTEGER NOT NULL,
			status TEXT CHECK (status IN ('delivered', 'received', 'read')) NOT NULL,
			type TEXT CHECK (type IN ('text', 'image')) NOT NULL,
			content TEXT NOT NULL,
			repliedTo INTEGER,
			forwardedFrom INTEGER,
			FOREIGN KEY (senderId) REFERENCES User(userId),
			FOREIGN KEY (conversationId) REFERENCES Conversation(conversationId),
			FOREIGN KEY (repliedTo) REFERENCES Message(messageId),
			FOREIGN KEY (forwardedFrom) REFERENCES Message(messageId));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating the message table: %w", err)
		}

		// Create Comment Table
		sqlStmt = `CREATE TABLE Comment (
			commentId INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp TEXT NOT NULL,
			senderId INTEGER NOT NULL, 
			messageId INTEGER NOT NULL,
			reaction TEXT NOT NULL,
			FOREIGN KEY (senderId) REFERENCES User(userId),
			FOREIGN KEY (messageId) REFERENCES Message(messageId));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating the comment table: %w", err)
		}

		// Create UserGroup Table
		sqlStmt = `CREATE TABLE UserGroup (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			userId INTEGER NOT NULL, 
			groupId INTEGER NOT NULL,
			FOREIGN KEY (userId) REFERENCES User(userId),
			FOREIGN KEY (groupId) REFERENCES GroupConversation(groupId));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating the UserGroup table: %w", err)
		}

		// Create MessageReceivers Table
		sqlStmt = `CREATE TABLE MessageReceivers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			userId INTEGER NOT NULL, 
			messageId INTEGER NOT NULL,
			FOREIGN KEY (userId) REFERENCES User(userId),
			FOREIGN KEY (messageId) REFERENCES Message(messageId));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating the MessageReceivers table: %w", err)
		}

		// Create MessageReaders Table
		sqlStmt = `CREATE TABLE MessageReaders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			userId INTEGER NOT NULL, 
			messageId INTEGER NOT NULL,
			FOREIGN KEY (userId) REFERENCES User(userId),
			FOREIGN KEY (messageId) REFERENCES Message(messageId));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating the MessageReaders table: %w", err)
		}
	}

	// Active the foreign keys
	sqlStmt := `PRAGMA foreign_keys = ON;`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("error activing the foreign keys: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
