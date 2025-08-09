# WASAText
![Go Version](https://img.shields.io/badge/Go-1.20%2B-blue)
![Vue.js Version](https://img.shields.io/badge/Vue.js-3.x-brightgreen)
![Docker](https://img.shields.io/badge/Docker-Ready-blue)

## Description 📖

Connect with your friends effortlessly using WASAText! Send and receive messages, whether one-on-one or in groups, all from the convenience of your PC. Enjoy seamless conversations with text or images and easily stay in touch through your private chats or group discussions!

## About This Project

This project’s initial structure comes from the [Fantastic Coffee (decaffeinated) template](https://github.com/sapienzaapps/fantastic-coffee-decaffeinated), adapted to fit WASAText’s requirements.

## Project Structure 🏗️

- **cmd/**  
  Houses all executable Go programs. These are entry points for the application and should focus only on "executable-stuff," such as:
  - Parsing CLI arguments
  - Reading environment variables
  - Initializing configuration before delegating logic to service packages  

- **doc/**  
  Contains project documentation and API specifications.  
  - Includes the `api.yaml` file, typically an api specification describing available endpoints and request/response formats.

- **service/**  
  Implements the backend logic and core functionalities of the project.  
  - Organized into sub-packages for different domains or features.  
  - Handles business rules, data processing, and integration with external systems.

- **vendor/**  
  Managed automatically by Go's dependency system.  
  - Stores a copy of all external dependencies for reproducible builds.  
  - Ensures the project compiles with consistent versions of libraries.

- **webui/**  
  Contains the implementation of the frontend user interface.  
  - Built using Vue.js  
  - May include components, pages, and assets such as stylesheets and icons  
  - Communicates with backend services via API calls

## Features 🌟

- **Conversation List** — Displays private or group chats in reverse chronological order with:
  - Username or group name  
  - Profile photo or group image  
  - Date and time of the latest message  
  - Text preview or photo icon

- **One-on-One & Group Chats**  
  - Start conversations with any WASAText user  
  - Create groups with any number of users  
  - Add members to groups (only existing members can invite)  
  - Leave groups at any time

- **Search** — Find users by username and view all existing WASAText usernames.

- **Messaging** —  
  - View all exchanged messages in reverse chronological order  
  - Messages include timestamps and sender usernames
  - Support for text and photo messages  
  - Add reactions/comments to messages, with visible names of users who reacted

## How To Run 🛠️

**Clone the Repository**
   ```bash
   git clone https://github.com/LucaMica02/WASAText.git
   ```
   ```bash
   cd WASAText
   ```
   
### Development Mode
1. **Backend**
  - Run:
   ```bash
   ./open-node.sh
   ```
3. **Frontend**
  - Build:
   ```bash
   ./open-node.sh
   ```
  - Run:
   ```bash
   yarn run dev
   ```
3. Go on http://localhost:5173/ and enjoy!

### Production Mode
1. **Backend**
  - Build:
   ```bash
   docker build -t backend:latest -f Dockerfile.backend .
   ```
  - Run:
   ```bash
   docker run -it --rm -p 3000:3000 backend:latest
   ```
2. **Frontend**
  - Build:
   ```bash
   docker build -t frontend:latest -f Dockerfile.frontend .
   ```
  - Run:
   ```bash
   docker run -it --rm -p 8080:80 frontend:latest
   ```
3. Go on http://localhost:8080/ and enjoy!

## Conclusion 🔚
**Enjoy building, exploring, and connecting with WASAText! 💬🚀**  
Stay in touch, collaborate with friends or teams, and make the most of the features this project offers.
