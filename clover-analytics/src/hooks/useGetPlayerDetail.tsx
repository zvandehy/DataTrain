import { gql, useQuery } from "@apollo/client";
// import Box from "@mui/material/Box";
// import CircularProgress from "@mui/material/CircularProgress";
// import moment from "moment";
// import { FilterGames } from "../shared/functions/filters.fn";
// import { CustomCalculation } from "../shared/interfaces/custom-prediction.interface";
// import { GameFilter } from "../shared/interfaces/graphql/filters.interface";
// import { Player } from "../shared/interfaces/graphql/player.interface";
// import { Projection } from "../shared/interfaces/graphql/proposition.interface";

// export const GET_PLAYER = gql`
//   query GetPlayer($playerID: Int!, $season: String!, $date: String!) {
//     player(input: { playerID: $playerID }) {
//       name
//       playerID
//       position
//       currentTeam {
//         abbreviation
//         teamID
//         name
//         injuries {
//           startDate
//           returnDate
//           status
//           player {
//             name
//           }
//         }
//       }
//       similarPlayers(input: { season: $season }) {
//         name
//         playerID
//         games(input: { season: $season, endDate: $date }) {
//           player {
//             name
//           }
//           points
//           season
//           assists
//           assist_percentage
//           rebounds
//           offensive_rebounds
//           offensive_rebound_percentage
//           defensive_rebounds
//           defensive_rebound_percentage
//           personal_fouls_drawn
//           steals
//           blocks
//           turnovers
//           opponent {
//             abbreviation
//             teamID
//           }
//           team {
//             abbreviation
//             teamID
//           }
//           minutes
//           date
//           field_goals_attempted
//           field_goal_percentage
//           field_goals_made
//           three_pointers_attempted
//           three_pointers_made
//           free_throws_attempted
//           free_throws_made
//           free_throws_percentage
//           usage
//           three_point_percentage
//           effective_field_goal_percentage
//           playoffs
//           home_or_away
//           # teamStats {
//           #   points
//           # }
//         }
//       }
//       projections(input: {}) {
//         date
//         opponent {
//           abbreviation
//           teamID
//           name
//           injuries {
//             startDate
//             returnDate
//             status
//             player {
//               name
//             }
//           }
//           similarTeams(input: { season: $season }) {
//             teamID
//             name
//             abbreviation
//           }
//         }
//         propositions {
//           target
//           type
//           sportsbook
//           lastModified
//         }
//         result {
//           player {
//             name
//           }
//           points
//           season
//           assists
//           assist_percentage
//           rebounds
//           offensive_rebounds
//           offensive_rebound_percentage
//           defensive_rebounds
//           defensive_rebound_percentage
//           personal_fouls_drawn
//           steals
//           blocks
//           turnovers
//           opponent {
//             abbreviation
//             teamID
//           }
//           team {
//             abbreviation
//             teamID
//           }
//           minutes
//           date
//           field_goals_attempted
//           field_goal_percentage
//           field_goals_made
//           three_pointers_attempted
//           three_pointers_made
//           free_throws_attempted
//           free_throws_made
//           free_throws_percentage
//           usage
//           three_point_percentage
//           effective_field_goal_percentage
//           playoffs
//           # teamStats {
//           #   points
//           # }
//         }
//       }
//       games(input: { season: $season }) {
//         player {
//           name
//         }
//         points
//         season
//         assists
//         assist_percentage
//         rebounds
//         offensive_rebounds
//         offensive_rebound_percentage
//         defensive_rebounds
//         defensive_rebound_percentage
//         personal_fouls_drawn
//         steals
//         blocks
//         turnovers
//         # teamStats {
//         #   points
//         #   assists
//         #   rebounds
//         #   three_pointers_attempted
//         #   blocks
//         #   steals
//         # }
//         team {
//           abbreviation
//           teamID
//           name
//         }
//         opponent {
//           abbreviation
//           teamID
//           name
//           injuries {
//             startDate
//             returnDate
//             status
//             player {
//               name
//             }
//           }
//         }
//         minutes
//         margin
//         home_or_away
//         date
//         field_goals_attempted
//         field_goal_percentage
//         field_goals_made
//         three_pointers_attempted
//         three_pointers_made
//         free_throws_attempted
//         free_throws_made
//         free_throws_percentage
//         usage
//         three_point_percentage
//         effective_field_goal_percentage
//         playoffs
//         # teamStats {
//         #   points
//         # }
//       }
//     }
//   }
// `;

// export interface QueryResult {
//   loading?: any;
//   error?: any;
//   data?: any;
// }

// export interface ProjectionQueryResult extends QueryResult {
//   data?: Player;
// }

// export const useGetPlayerDetails = ({
//   playerID,
//   predictionFilter,
//   gameFilter,
//   customModel,
// }: {
//   playerID: number;
//   predictionFilter: GameFilter;
//   gameFilter: GameFilter;
//   customModel: CustomCalculation;
// }): ProjectionQueryResult => {
//   const { loading, error, data } = useQuery(GET_PLAYER, {
//     //TODO: refactor query to use predictionFilter
//     variables: {
//       playerID: playerID,
//       // season: predictionFilter.season,
//       date: predictionFilter.endDate,
//     },
//   });
//   let loadingComponent;
//   if (loading) {
//     loadingComponent = (
//       <Box className={"loading-results"}>
//         <h1>Loading </h1>
//         <CircularProgress />
//       </Box>
//     );
//   }
//   let errorComponent;
//   if (error) {
//     console.error(JSON.stringify(error) + error.message);
//     errorComponent = <Box>{JSON.stringify(error) + error.message}</Box>;
//   }
//   if (data && data?.player?.projections) {
//     let projections: Projection[] = data.player.projections.map(
//       (projection: Projection) => {
//         let p = {
//           ...data.player,
//           games: FilterGames(data.player.games, {
//             endDate: projection.date,
//           }).sort((a, b) => {
//             return moment(a.date).diff(b.date);
//           }),
//         };
//         let newProjection: Projection = {
//           ...projection,
//           player: p,
//         };
//         return newProjection;
//       }
//     );
//     const player: Player = {
//       ...data.player,
//       projections: projections,
//     };
//     return {
//       loading: loadingComponent,
//       error: errorComponent,
//       data: player,
//     };
//   }
//   return {
//     loading: loadingComponent,
//     error: errorComponent,
//   };
// };
