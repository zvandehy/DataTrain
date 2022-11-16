import { Label } from "@material-ui/icons";
import {
  Dialog,
  DialogTitle,
  ToggleButtonGroup,
  ToggleButton,
  TextField,
  DialogContentText,
  DialogActions,
  Button,
  Typography,
  Box,
} from "@mui/material";
import { useReducer } from "react";
import { DEFAULT_MODEL } from "../../shared/constants";
import {
  CustomCalculation,
  Factor,
  ModelInput,
} from "../../shared/interfaces/custom-prediction.interface";
import "./custom-model-dialog.component.css";
import {
  CustomModelActionType,
  customModelReducer,
  INITIAL_CUSTOM_MODEL_STATE,
} from "./custom-model-dialog.reducer";

interface CustomModelDialogProps {
  open: boolean;
  closeDialog: () => void;
  setCustomModel: (customModel: ModelInput) => void;
}

const CustomModelDialog: React.FC<CustomModelDialogProps> = ({
  open,
  closeDialog,
  setCustomModel,
}: CustomModelDialogProps) => {
  // const [customModelForm, dispatch] = useReducer(
  //   customModelReducer,
  //   localStorage.getObject("customModel") ?? DEFAULT_MODEL
  // );
  let customModelForm = DEFAULT_MODEL;
  const dispatchLog = (action: any) => {
    // dispatch(action);
    console.log(action);
  };
  const handleClose = (event: any, reason: string) => {
    if (reason === "backdropClick") {
      return;
    }
    closeDialog();
  };

  const cancel = () => {
    dispatchLog({ type: CustomModelActionType.RESET });
    closeDialog();
  };

  const save = () => {
    setCustomModel(customModelForm);
    closeDialog();
  };

  const customModelTotalWeight = (customModel: CustomCalculation) => {
    let total = 0;
    customModel.recency?.forEach((r) => {
      total += r.weight;
    });
    customModel.recencyPct?.forEach((r) => {
      total += r.weight;
    });
    total += customModel.homeAwayWeight ?? 0;
    total += customModel.opponentWeight ?? 0;
    total += customModel.similarPlayers?.weight ?? 0;
    total += customModel.similarTeams?.weight ?? 0;
    return total;
  };

  return (
    <Dialog
      fullWidth={true}
      open={open}
      onClose={handleClose}
      sx={{
        // "& .MuiPaper-root": { backgroundColor: "gray" },
        color: "black",
      }}
    >
      <DialogTitle color="black">Custom Model Calculation</DialogTitle>
      <Typography color="black" textAlign={"left"}>
        Include Push?
      </Typography>
      {/* <ToggleButtonGroup
        color="primary"
        value={customModelForm.includePush.toString()}
        exclusive
        onChange={(e, value: boolean) => {
          dispatchLog({
            type: CustomModelActionType.TOGGLE_INCLUDE_PUSH,
            payload: value,
          });
        }}
      >
        <ToggleButton value="true">Yes</ToggleButton>
        <ToggleButton value="false">No</ToggleButton>
      </ToggleButtonGroup> */}
      {/* <Typography color="black" textAlign="left">
        Include Games on Different Teams?
      </Typography>
      <ToggleButtonGroup
        color="primary"
        value={customModelForm.includeOnDifferentTeam.toString()}
        exclusive
        onChange={(e, value: boolean) => {
          dispatchLog({
            type: CustomModelActionType.TOGGLE_INCLUDE_DIFF_TEAM,
            payload: value,
          });
        }}
      >
        <ToggleButton value="true">Yes</ToggleButton>
        <ToggleButton value="false">No</ToggleButton>
      </ToggleButtonGroup> */}
      <Typography color="black" textAlign="left">
        Last X Games Breakdown
      </Typography>
      {/* <Box>
        {customModelForm.recency?.length ? (
          customModelForm.recency.map((fragment: Factor, index: number) => {
            return (
              <Box key={"Recency: " + index}>
                <TextField
                  id={`recency-count-input-${index}`}
                  label="Last X Games"
                  type="number"
                  onChange={(e) => {
                    let newRecency = customModelForm.recency!;
                    newRecency[index].count = parseInt(e.target.value);
                    dispatchLog({
                      type: CustomModelActionType.SET_RECENCY,
                      payload: newRecency,
                    });
                  }}
                  value={fragment.count ?? 0}
                  sx={{
                    "& .MuiInputBase-input": { color: "black" },
                  }}
                  variant={"filled"}
                />
                <TextField
                  id={`recency-weight-input-${index}`}
                  label="Weight"
                  type="number"
                  onChange={(e) => {
                    let newRecency = customModelForm.recency!;
                    newRecency[index].weight = parseInt(e.target.value);
                    dispatchLog({
                      type: CustomModelActionType.SET_RECENCY,
                      payload: newRecency,
                    });
                  }}
                  value={fragment.weight ?? 0}
                  sx={{
                    "& .MuiInputBase-input": { color: "black" },
                  }}
                  variant={"filled"}
                />
              </Box>
            );
          })
        ) : (
          <></>
        )}
      </Box> */}
      <Typography color="black" textAlign="left">
        Players with Similar Statistics vs. Opponent
      </Typography>
      {/* <TextField
        id="similar-players-weight-input"
        label="Similar Players Weight"
        type="number"
        onChange={(e) => {
          dispatchLog({
            type: CustomModelActionType.SET_SIMILAR_PLAYERS_WEIGHT,
            payload: parseInt(e.target.value),
          });
        }}
        value={customModelForm.similarPlayers?.weight ?? 0}
        sx={{
          "& .MuiInputBase-input": { color: "black" },
        }}
        variant={"filled"}
      /> */}
      <Typography color="black" textAlign="left">
        Player vs. Teams Similar to the Opponent
      </Typography>
      {/* <TextField
        id="similar-teams-weight-input"
        label="Similar Teams Weight"
        type="number"
        onChange={(e) => {
          dispatchLog({
            type: CustomModelActionType.SET_SIMILAR_TEAMS_WEIGHT,
            payload: parseInt(e.target.value),
          });
        }}
        value={customModelForm.similarTeams?.weight ?? 0}
        sx={{
          "& .MuiInputBase-input": { color: "black" },
        }}
        variant={"filled"}
      /> */}
      <Typography color="black" textAlign="left">
        Previous Head to Head Matchups
      </Typography>
      {/* <TextField
        id="versus-opponent-weight-input"
        label="Versus Opponent Weight"
        type="number"
        onChange={(e) => {
          dispatchLog({
            type: CustomModelActionType.SET_VERSUS_OPPONENT_WEIGHT,
            payload: parseInt(e.target.value),
          });
        }}
        value={customModelForm.opponentWeight ?? 0}
        sx={{
          "& .MuiInputBase-input": { color: "black" },
        }}
        variant={"filled"}
      /> */}
      {/* <DialogContentText>
        Total Weight (must equal 100%):{" "}
        {customModelTotalWeight(customModelForm)}
      </DialogContentText> */}
      {/* <DialogActions>
        {customModelTotalWeight(customModelForm) !== 100 ? (
          <Button variant={"contained"} disabled>
            Save
          </Button>
        ) : (
          <Button variant={"contained"} onClick={save}>
            Save
          </Button>
        )}
        <Button variant={"outlined"} onClick={cancel}>
          Cancel
        </Button>
      </DialogActions> */}
    </Dialog>
  );
};

export default CustomModelDialog;
