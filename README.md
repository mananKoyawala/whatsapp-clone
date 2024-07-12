# Whatsapp-Clone

- The presenting backend is the conceptial implementation of chatting functionality like Whatsapp. It's utilizies the websockets protocol.

## Language, Tools and Services are used

**1. Go (Golang) :** I used golang to develop the my backend which is way power fully handle all the things like Connectivity, High performance, Concurrency, Usability and Simplicity.

**2. Postgres (database) :** I used Postgres which is RDBMS, helps to use my SQL knowledge in it.

**3. Docker :** I used docker to get Postgres Database as service. So we don't need to install the postgres in our local system

**4. Database Migration :** Golang Database migration is very useful when we need to add new table or even we have to add new column , remove or change something in column. It's make it very easy.

**5. AWS :** I used AWS s3 bucket to store all the files that are sharing this system to keep store at one secure place. As of now the system can store only .png, .jpg, .jepg files.

**6. Logging :** Used log/slog (STLB) to implement logging to further analyse it.

- It has 4 levels

1.  **Error** :- Logged when error occurs
2.  **Warn** :- Warned when something was expected but not fullfiled like validations
3.  **Info** :- Logged when the everything was good
4.  **Debug** :- Logged when some infomation are retrived at repository level

## Learnings

- During developing this backend, I learned couple of tools like Docker, Database Migration, Logging, AWS (s3 bucket), Websocket, How to use Self-Signed certificate (HTTPs) and Data encryption using AES.
- I used the Clean-code architecture, means every things are divied into layer
- There are three layers

1. Repository layer :- Deals with only database.
2. Service layer :- Act as intermediate layer between Repository and Handler layer.
3. Handler layer :- Handle the request that upcoming from client. I used _GIN_ for that.

## Documentation

For more details, please read the [documentation](https://manankoyawala.github.io/whastapp-clone-doc).

In the documentation, you will find a detailed overview of what was done in the project and the various components included.

## License

[MIT License](LICENSE)
