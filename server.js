const express = require('express')
const path = require('path');
var cors = require('cors')

const app = express()
const port = process.env.PORT || 3000

app.use(express.static("public"));

app.use(cors())
app.listen(port, () => {
    console.log(`Example app listening at http://localhost:${port}`)
})