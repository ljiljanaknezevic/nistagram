const mongoose = require("mongoose");
const User = require("../model/User.model");
const connection = "mongodb://mongo:27017/service1";
const connectDb = () => {
 return mongoose.connect(connection, { useNewUrlParser: true });
};
module.exports = connectDb;