import { CustomCalculation } from "../../shared/interfaces/custom-prediction.interface";
export const INITIAL_CUSTOM_MODEL_STATE = {
  includePush: true,
  includeOnDifferentTeam: true,
  recency: [
    { count: 0, weight: 20 },
    { count: 20, weight: 10 },
    { count: 10, weight: 10 },
    { count: 5, weight: 12 },
  ],
  similarPlayers: { count: 10, weight: 12 },
  similarTeams: { count: 3, weight: 15 },
  opponentWeight: 21,
  // homeAwayWeight:0.1,
};

export enum CustomModelActionType {
  SET_SIMILAR_PLAYERS_WEIGHT = "SET_SIMILAR_PLAYERS_WEIGHT",
  SET_SIMILAR_PLAYERS_COUNT = "SET_SIMILAR_PLAYERS_COUNT",
  SET_SIMILAR_TEAMS_WEIGHT = "SET_SIMILAR_TEAMS_WEIGHT",
  SET_SIMILAR_TEAMS_COUNT = "SET_SIMILAR_TEAMS_COUNT",
  SET_VERSUS_OPPONENT_WEIGHT = "SET_VERSUS_OPPONENT_WEIGHT",
  TOGGLE_INCLUDE_PUSH = "TOGGLE_INCLUDE_PUSH",
  TOGGLE_INCLUDE_DIFF_TEAM = "TOGGLE_INCLUDE_DIFF_TEAM",
  SET_RECENCY_COUNT = "SET_RECENCY_COUNT",
  SET_RECENCY_WEIGHT = "SET_RECENCY_WEIGHT",
  SET_RECENCY = "SET_RECENCY",
  RESET = "RESET",
}

export interface CustomModelAction {
  type: CustomModelActionType;
  payload: any;
}

// export interface CustomCalculation {
//     includePush: boolean;
//     includeOnDifferentTeam: boolean;
//     recency?: Factor[];
//     recencyPct?: Factor[];
//     similarPlayers?: Factor;
//     similarTeams?: Factor;
//     homeAwayWeight?: number;
//     opponentWeight?: number;
//     playoffs?: Factor;
//   }

//   export interface Factor {
//     weight: number;
//     count?: number;
//     seasons?: string[];
//   }

export const customModelReducer = (
  state: CustomCalculation,
  action: CustomModelAction
): CustomCalculation => {
  switch (action.type) {
    case CustomModelActionType.SET_SIMILAR_PLAYERS_WEIGHT:
      return {
        ...state,
        similarPlayers: { ...state.similarPlayers, weight: action.payload },
      };
    case CustomModelActionType.SET_SIMILAR_TEAMS_WEIGHT:
      return {
        ...state,
        similarTeams: { ...state.similarTeams, weight: action.payload },
      };
    case CustomModelActionType.SET_SIMILAR_PLAYERS_COUNT:
      return {
        ...state,
        similarPlayers: {
          ...state.similarPlayers,
          weight: state.similarPlayers?.weight ?? 0,
          count: action.payload,
        },
      };
    case CustomModelActionType.SET_SIMILAR_TEAMS_COUNT:
      return {
        ...state,
        similarTeams: {
          ...state.similarTeams,
          weight: state.similarTeams?.weight ?? 0,
          count: action.payload,
        },
      };
    case CustomModelActionType.SET_VERSUS_OPPONENT_WEIGHT:
      return {
        ...state,
        opponentWeight: action.payload,
      };
    case CustomModelActionType.TOGGLE_INCLUDE_PUSH:
      return {
        ...state,
        includePush: action.payload === "true",
      };
    case CustomModelActionType.TOGGLE_INCLUDE_DIFF_TEAM:
      return {
        ...state,
        includeOnDifferentTeam: action.payload === "true",
      };
    case CustomModelActionType.SET_RECENCY:
      return {
        ...state,
        recency: action.payload,
      };
    case CustomModelActionType.RESET:
      return (
        localStorage.getObject("customModel") ?? INITIAL_CUSTOM_MODEL_STATE
      );
  }
  return state;
};
