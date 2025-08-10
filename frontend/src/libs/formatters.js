import { sprintf } from "sprintf-js";

/**
 * Formats engine hours to either hh.h or hh:mm depending
 * on the format input specified.
 * @param {*} value the engine hour in hh.h format (e.g. 9.7)
 * @param {*} format either hh.h or hh:mm
 */
export const formatEngineHour = (value, format) => {
  if (value == null) {
    return null;
  }

  if (format == null || format == "hh.h") {
    return sprintf("%.1f", Number(value));
  } else if (format == "hh:mm") {
    const hours = parseInt(value);
    const minutes = (parseFloat(value) - hours) * 60;
    return hours + ":" + sprintf("%02d", Math.round(minutes));
  }
};

/**
 * Returns true if the provided value can be parsed into the provided format.
 * @param {*} value the engine hour
 * @param {*} format the format (e.g. hh.h, hh:mm)
 */
export const isValidEngineHour = (value, format) => {
  if (format == null || format == "hh.h") {
    if (value.match(/^\d+(?:\.\d+)?$/)) {
      return true;
    }
    return false;
  } else if (format == "hh:mm") {
    if (value.match(/^\d+:[0-5][0-9]$/)) {
      return true;
    }
    return false;
  }
  return false;
};

/**
 * Given a valid engine hour format, it returns our
 * default format 'hh.h'
 * @param {*} value an engine hour format
 */
export const convertEngineHour = (value) => {
  // convert hh:mm
  if (value.match(/^\d+:[0-5][0-9]$/)) {
    const hourMin = value.split(":");
    const hours = Number(hourMin[0]);
    const hoursDec = Number(hourMin[1]) / 60;
    return Number(hours + hoursDec);
  }
  if (value.match(/^\d+(?:\.\d+)?$/)) {
    return Number(value);
  }
  return null;
};

/**
 * Returns a representation of the engine hour format to be used as a label
 * for input fields.
 * @param {*} format hh.h or hh:mm
 */
export const formatEngineHourLabel = (format) => {
  if (format == null || format == "hh.h") {
    return "hrs";
  } else if (format == "hh:mm") {
    return format;
  }
  return "hrs";
};

/**
 * Formats a value of a given currency into a string with the currency appended.
 * @param {*} value the value
 * @param {*} currency the currency (e.g. CHF)
 */
export const formatCurrency = (value, currency) => {
  if (currency != null) {
    return sprintf("%.2f %s", Math.round(Number(value) * 100) / 100, currency);
  } else {
    return sprintf("%.2f", Math.round(Number(value) * 100) / 100);
  }
};

/**
 * Formats a value to be displayed as an amount of fuel in liters
 * @param {*} value fuel in liters
 */
export const formatFuel = (value) => {
  return sprintf("%.2f", Number(value));
};

/**
 * Formats a value to be displayed as liters / hour consumption
 * @param {*} value liters per hour
 */
export const formatFuelConsumption = (value) => {
  return sprintf("%.1f", Number(value));
};

/**
 * Formats a duration in seconds into a string of
 * "A days B hours C minutes D seconds".
 * @param {*} value seconds
 */
export const formatDurationString = (value, displaySeconds) => {
  // seconds
  const fmtSeconds = value % 60;
  // minutes
  let minutes = (value - fmtSeconds) / 60;
  const fmtMinutes = minutes % 60;
  // hours
  minutes = minutes - fmtMinutes;
  const fmtHours = (minutes / 60) % 24;
  // days
  minutes = minutes - fmtHours * 60;
  const fmtDays = minutes / 60 / 24;
  // return value
  switch (true) {
    case displaySeconds == true && fmtDays > 0:
      return sprintf(
        "%d days %d hours %d minutes %d seconds",
        fmtDays,
        fmtHours,
        fmtMinutes,
        fmtSeconds,
      );
    case displaySeconds == true:
      return sprintf(
        "%d hours %d minutes %d seconds",
        fmtHours,
        fmtMinutes,
        fmtSeconds,
      );
    case fmtDays > 0:
      return sprintf(
        "%d days %d hours %d minutes",
        fmtDays,
        fmtHours,
        fmtMinutes,
      );
    default:
      return sprintf("%d hours %d minutes", fmtHours, fmtMinutes);
  }
};

/**
 * Formats a number into the Locale representation.
 * @param {} number
 * @returns
 */
export const formatNumber = (number) => {
  const c = Number(number);
  // TODO: make this configurable.
  return c.toLocaleString("de-CH");
};

/**
 * Formats a number into number followed by the currency. It
 * uses formatNumber() internally.
 * @param {*} value
 * @param {*} currency
 * @returns
 */
export const formatCost = (value, currency) => {
  const v = parseFloat(value);
  if (isNaN(v)) {
    return "NaN" + " " + currency;
  }
  const n = formatNumber(v);
  // TODO make this configurable.
  return sprintf("%.02f %s", n, currency);
};
