const express = require('express');

const connectDb = require('./src/utils/connection');
const User = require('./src/model/User.model');

const app = express();

app.get('/', (req, res) => {
  res.send('hello from service 1');
});

app.get('/service1', async (req, res) => {
  const users = await User.find();
  res.json(users);
});

app.listen(3000, () => {
  console.log('Listening on 3000');
  connectDb()
    .then(() => {
      console.log('MongoDb connected');
    })
});
