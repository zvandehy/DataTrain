import { CheckBox } from "@material-ui/icons";
import DisabledByDefaultIcon from "@mui/icons-material/DisabledByDefault";
import UnfoldLessDoubleIcon from "@mui/icons-material/UnfoldLessDouble";
import CrisisAlertIcon from "@mui/icons-material/CrisisAlert";
import { Icon, Typography } from "@mui/material";
import { makeStyles } from "@mui/styles";
import { COLORS } from "../../shared/styles/constants";

const hitMissStyle = makeStyles({
  hit: {
    color: COLORS.HIGHER,
  },
  miss: {
    color: COLORS.LOWER,
  },
  push: {
    color: COLORS.PUSH,
    transform: "rotate(90deg)",
  },
});

interface HitMissIconProps {
  outcome: string;
  result?: number;
}

export const HitMissIcon: React.FC<HitMissIconProps> = ({
  outcome,
}: HitMissIconProps) => {
  const classes = hitMissStyle();
  return (
    <Icon sx={{ m: "auto" }}>
      {outcome.toLowerCase() === "hit" ? (
        <CheckBox className={classes.hit} />
      ) : outcome.toLowerCase() === "miss" ? (
        <DisabledByDefaultIcon className={classes.miss} />
      ) : outcome.toLowerCase() === "push" ? (
        <UnfoldLessDoubleIcon className={classes.push} />
      ) : (
        <> </>
      )}
    </Icon>
  );
};
// TODO: Return an icon instead of a react element
export const HitMissTypography: React.FC<HitMissIconProps> = ({
  outcome,
  result,
}: HitMissIconProps) => {
  return (
    <Typography
      variant={"overline"}
      fontSize={14}
      fontWeight={800}
      paddingRight={1}
      textTransform={"uppercase"}
    >
      {outcome} {result !== undefined ? "(" + result.toFixed(2) + ")" : ""}
    </Typography>
  );
};
