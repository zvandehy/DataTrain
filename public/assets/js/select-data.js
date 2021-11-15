(function () {
    "use strict"; // Start of use strict

    $("#player-select").selectize({
        valueField: "id",
        searchField: ["first_name", "last_name"],
        searchConjuntion: "or",
        create: false,
        options: [],
        render: {
            option: function (item, escape) {
                return (
                    "<div>" + escape(item.first_name) + " " + escape(item.last_name) + "</div>"
                );
            },
            item: function (selected, escape) {
                return `<span>${escape(selected.first_name)} ${escape(selected.last_name)}</span>`
            }

        },
        load: function (query, callback) {
            console.log("query: " + JSON.stringify(query));
            if (!query.length) return callback();
            $.ajax({
                url: `http://stats.nba.com/stats/scoreboard/?GameDate=02/14/2015&LeagueID=00&DayOffset=0`,
                type: "GET",
                // crossDomain: true,
                headers: {
                    "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:93.0) Gecko/20100101 Firefox/93.0",
                    "Referer": "https://stats.nba.com/",

                },
                error: function (err) {
                    console.log("Error loading players");
                    console.log(err)
                    callback();
                },
                success: function (res) {
                    console.log(res)
                    callback();
                    // let data = res["included"]
                    // let results = []
                    // for (let i = 0; i < data; i++) {
                    //     data[i] = item
                    //     if (item.type == "new_player") {
                    //         results.push({ "id": item.id, "position": item.attributes.position, "name": item.attributes.name, "team": item.attributes.team })
                    //     }
                    // }

                    // callback(results);
                },
            });
        },
    });

    $("#stat-select").selectize({
        create: true,
        sortField: {
            field: 'text',
            direction: 'asc'
        },
        dropdownParent: 'body'
    });

    $("#season-select").selectize({
        create: true,
        dropdownParent: 'body'
    });

    $("#opponent-select").selectize({
        valueField: "abbreviation",
        labelField: "name",
        searchField: ["abbreviation", "city", "full_name", "name"],
        searchConjuntion: "or",
        create: false,
        options: teams,
        items: teams,
        closeAfterSelection: true,
        render: {
            option: function (item, escape) {
                return (
                    "<div>" + escape(item.full_name) + "</div>"
                );
            },

        },
        load: function (query, callback) {
            callback(teams)
        },
        preload: true,
    });
})();

function analyze() {
    let opponent = $("#opponent-select")[0].firstChild.value
    let player = $("#player-select")[0].firstChild.value
    let stat = $("#stat-select")[0].firstChild.value
    let season = $("#season-select")[0].firstChild.value
    localStorage.setItem("opponent", JSON.stringify(opponent));
    localStorage.setItem("player", JSON.stringify(player));
    localStorage.setItem("stat", JSON.stringify(stat));
    localStorage.setItem("season", JSON.stringify(season));

    let playerSpans = $(".player-name")
    for (let i = 0; i < playerSpans.length; i++) {
        playerSpans[i].innerHTML = JSON.parse(localStorage.getItem("player"));
    }
    let statSpans = $(".stat")
    for (let i = 0; i < statSpans.length; i++) {
        statSpans[i].innerHTML = JSON.parse(localStorage.getItem("stat"));
    }
    let opponentSpans = $(".opponent")
    for (let i = 0; i < opponentSpans.length; i++) {
        opponentSpans[i].innerHTML = JSON.parse(localStorage.getItem("opponent"));
    }
    let gamelogs = [{ "date": "10/23/2021", "assists": 7 }, { "date": "10/24/2021", "assists": 8 }, { "date": "10/25/2021", "assists": 9 }, { "date": "10/26/2021", "assists": 10 }]
    localStorage.setItem("gamelogs", JSON.stringify(gamelogs));
}

function storageAvailable(type) {
    var storage;
    try {
        storage = window[type];
        var x = '__storage_test__';
        storage.setItem(x, x);
        storage.removeItem(x);
        return true;
    }
    catch (e) {
        return e instanceof DOMException && (
            // everything except Firefox
            e.code === 22 ||
            // Firefox
            e.code === 1014 ||
            // test name field too, because code might not be present
            // everything except Firefox
            e.name === 'QuotaExceededError' ||
            // Firefox
            e.name === 'NS_ERROR_DOM_QUOTA_REACHED') &&
            // acknowledge QuotaExceededError only if there's something already stored
            (storage && storage.length !== 0);
    }
}

function localStorageSpace() {
    var allStrings = '';
    for (var key in window.localStorage) {
        if (window.localStorage.hasOwnProperty(key)) {
            allStrings += window.localStorage[key];
        }
    }
    return allStrings ? 3 + ((allStrings.length * 16) / (8 * 1024)) + ' KB' : 'Empty (0 KB)';
};


var teams = [
    {
        "id": 1,
        "abbreviation": "ATL",
        "city": "Atlanta",
        "conference": "East",
        "division": "Southeast",
        "full_name": "Atlanta Hawks",
        "name": "Hawks"
    },
    {
        "id": 2,
        "abbreviation": "BOS",
        "city": "Boston",
        "conference": "East",
        "division": "Atlantic",
        "full_name": "Boston Celtics",
        "name": "Celtics"
    },
    {
        "id": 3,
        "abbreviation": "BKN",
        "city": "Brooklyn",
        "conference": "East",
        "division": "Atlantic",
        "full_name": "Brooklyn Nets",
        "name": "Nets"
    },
    {
        "id": 4,
        "abbreviation": "CHA",
        "city": "Charlotte",
        "conference": "East",
        "division": "Southeast",
        "full_name": "Charlotte Hornets",
        "name": "Hornets"
    },
    {
        "id": 5,
        "abbreviation": "CHI",
        "city": "Chicago",
        "conference": "East",
        "division": "Central",
        "full_name": "Chicago Bulls",
        "name": "Bulls"
    },
    {
        "id": 6,
        "abbreviation": "CLE",
        "city": "Cleveland",
        "conference": "East",
        "division": "Central",
        "full_name": "Cleveland Cavaliers",
        "name": "Cavaliers"
    },
    {
        "id": 7,
        "abbreviation": "DAL",
        "city": "Dallas",
        "conference": "West",
        "division": "Southwest",
        "full_name": "Dallas Mavericks",
        "name": "Mavericks"
    },
    {
        "id": 8,
        "abbreviation": "DEN",
        "city": "Denver",
        "conference": "West",
        "division": "Northwest",
        "full_name": "Denver Nuggets",
        "name": "Nuggets"
    },
    {
        "id": 9,
        "abbreviation": "DET",
        "city": "Detroit",
        "conference": "East",
        "division": "Central",
        "full_name": "Detroit Pistons",
        "name": "Pistons"
    },
    {
        "id": 10,
        "abbreviation": "GSW",
        "city": "Golden State",
        "conference": "West",
        "division": "Pacific",
        "full_name": "Golden State Warriors",
        "name": "Warriors"
    },
    {
        "id": 11,
        "abbreviation": "HOU",
        "city": "Houston",
        "conference": "West",
        "division": "Southwest",
        "full_name": "Houston Rockets",
        "name": "Rockets"
    },
    {
        "id": 12,
        "abbreviation": "IND",
        "city": "Indiana",
        "conference": "East",
        "division": "Central",
        "full_name": "Indiana Pacers",
        "name": "Pacers"
    },
    {
        "id": 13,
        "abbreviation": "LAC",
        "city": "LA",
        "conference": "West",
        "division": "Pacific",
        "full_name": "LA Clippers",
        "name": "Clippers"
    },
    {
        "id": 14,
        "abbreviation": "LAL",
        "city": "Los Angeles",
        "conference": "West",
        "division": "Pacific",
        "full_name": "Los Angeles Lakers",
        "name": "Lakers"
    },
    {
        "id": 15,
        "abbreviation": "MEM",
        "city": "Memphis",
        "conference": "West",
        "division": "Southwest",
        "full_name": "Memphis Grizzlies",
        "name": "Grizzlies"
    },
    {
        "id": 16,
        "abbreviation": "MIA",
        "city": "Miami",
        "conference": "East",
        "division": "Southeast",
        "full_name": "Miami Heat",
        "name": "Heat"
    },
    {
        "id": 17,
        "abbreviation": "MIL",
        "city": "Milwaukee",
        "conference": "East",
        "division": "Central",
        "full_name": "Milwaukee Bucks",
        "name": "Bucks"
    },
    {
        "id": 18,
        "abbreviation": "MIN",
        "city": "Minnesota",
        "conference": "West",
        "division": "Northwest",
        "full_name": "Minnesota Timberwolves",
        "name": "Timberwolves"
    },
    {
        "id": 19,
        "abbreviation": "NOP",
        "city": "New Orleans",
        "conference": "West",
        "division": "Southwest",
        "full_name": "New Orleans Pelicans",
        "name": "Pelicans"
    },
    {
        "id": 20,
        "abbreviation": "NYK",
        "city": "New York",
        "conference": "East",
        "division": "Atlantic",
        "full_name": "New York Knicks",
        "name": "Knicks"
    },
    {
        "id": 21,
        "abbreviation": "OKC",
        "city": "Oklahoma City",
        "conference": "West",
        "division": "Northwest",
        "full_name": "Oklahoma City Thunder",
        "name": "Thunder"
    },
    {
        "id": 22,
        "abbreviation": "ORL",
        "city": "Orlando",
        "conference": "East",
        "division": "Southeast",
        "full_name": "Orlando Magic",
        "name": "Magic"
    },
    {
        "id": 23,
        "abbreviation": "PHI",
        "city": "Philadelphia",
        "conference": "East",
        "division": "Atlantic",
        "full_name": "Philadelphia 76ers",
        "name": "76ers"
    },
    {
        "id": 24,
        "abbreviation": "PHX",
        "city": "Phoenix",
        "conference": "West",
        "division": "Pacific",
        "full_name": "Phoenix Suns",
        "name": "Suns"
    },
    {
        "id": 25,
        "abbreviation": "POR",
        "city": "Portland",
        "conference": "West",
        "division": "Northwest",
        "full_name": "Portland Trail Blazers",
        "name": "Trail Blazers"
    },
    {
        "id": 26,
        "abbreviation": "SAC",
        "city": "Sacramento",
        "conference": "West",
        "division": "Pacific",
        "full_name": "Sacramento Kings",
        "name": "Kings"
    },
    {
        "id": 27,
        "abbreviation": "SAS",
        "city": "San Antonio",
        "conference": "West",
        "division": "Southwest",
        "full_name": "San Antonio Spurs",
        "name": "Spurs"
    },
    {
        "id": 28,
        "abbreviation": "TOR",
        "city": "Toronto",
        "conference": "East",
        "division": "Atlantic",
        "full_name": "Toronto Raptors",
        "name": "Raptors"
    },
    {
        "id": 29,
        "abbreviation": "UTA",
        "city": "Utah",
        "conference": "West",
        "division": "Northwest",
        "full_name": "Utah Jazz",
        "name": "Jazz"
    },
    {
        "id": 30,
        "abbreviation": "WAS",
        "city": "Washington",
        "conference": "East",
        "division": "Southeast",
        "full_name": "Washington Wizards",
        "name": "Wizards"
    }
]