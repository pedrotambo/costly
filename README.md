![Build](https://github.com/pedrotambo/costly/actions/workflows/go.yml/badge.svg)


# Costly

Welcome to the Costly project! This software aims to assist restaurants in calculating the costs of their recipes to optimize pricing and manage expenses effectively. This project is still in progress and might evolve in a bigger system with more functionalities or become part of multiple micro-services aimed to help restaurant management. Some aspects might be considered overdesigned because I'm using this project to experiment and put into practice new designs and patterns.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Technologies Used](#technologies-used)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Running the Application](#running-the-application)
- [Contributing](#contributing)
- [License](#license)

## Introduction

This project is designed to simplify the process of calculating recipe costs for restaurants. By utilizing this tool, restaurant owners and chefs can input the ingredients and quantities used in their recipes to obtain an accurate estimate of the total cost per serving. This information can be invaluable for setting menu prices, tracking expenses, and optimizing profitability.

## Features

- **Recipe Management:** Add, edit, and delete recipes effortlessly.
- **Ingredient Tracking:** Easily manage and update ingredient information.
- **Cost Calculation:** Automatically calculate the total cost of each recipe based on ingredient prices.

## Technologies Used

- **Backend:** Go (Golang), SQLite
- **Frontend:** React, TypeScript

## Getting Started

Follow these instructions to get a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go installed on your system. You can download it from [here](https://golang.org/doc/install).
- Node.js and npm installed on your system. You can download them from [here](https://nodejs.org/).
- Git installed on your system. You can download it from [here](https://git-scm.com/).
- Air (Live reload for Go apps). You can install it using instructions [here](https://github.com/cosmtrek/air?tab=readme-ov-file#installation).

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/pedrotambo/costly.git
   ```

2. Navigate to the project directory:

   ```bash
   cd costly
   ```

3. Install backend dependencies:

   ```bash
   go mod tidy
   ```

4. Install frontend dependencies:

   ```bash
   cd front
   npm install
   ```

### Running the Application

1. Start the backend server:

   ```bash
   air
   ```

2. Start the frontend development server:

   ```bash
   cd frontend
   npm start run
   ```

3. Access the application in your web browser at `http://localhost:3000`.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This project is licensed under the [MIT License](LICENSE).
