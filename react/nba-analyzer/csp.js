module.exports = {
  dev: {
  "default-src": ["'self'"],
  "style-src": [
    "'self'",
    "https://*.google.com",
  ]
  },
  prod: {
  "default-src": "'self'",  // can be either a string or an array.
  "style-src": [
    "'self'",
    "https://unpkg.com",
  ],
  "connect-src": [
    "'self'",
    "https://datatrain-nba-yxh2z.ondigitalocean.app"
  ],
  "script-src": [
    "'self'",
  ],
  "img-src": [
    "'self'",
    "ak-static.cms.nba.com",
  ],
  "object-src": "'none'"
  }
}