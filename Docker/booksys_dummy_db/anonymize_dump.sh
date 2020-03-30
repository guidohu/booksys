#!/bin/bash

OPTIND=1         # Reset in case getopts has been used previously in the shell.

# Initialize our own variables:
input_file=""
output_file=""
verbose=0

show_help() {
    cat << 'EOF'
    usage: $0 [-hv] -i <input_file> -o <output_file>

    OPTIONS:
    -h          display this help
    -v          run in verbose mode
    -i          specify input file (path)
    -o          specify output file (path)
EOF
}

anonymize() {
    # create two empty files
    echo "" > $output_file
    tempfile="${output_file}.tmp"
    echo "" > $tempfile

    # initialize tracking variable
    deleting_browser_session_data=0
    deleting_password_reset_data=0

    # delete sensitive content
    while IFS= read -r line
    do  
        if [[ "$line" =~ 'INSERT INTO `browser_session`' ]]; then
            [[ $verbose = 1 ]] && echo "START DELETION:"
            deleting_browser_session_data=1
        fi
        if [[ "$line" =~ 'INSERT INTO `password_reset`' ]]; then
            [[ $verbose = 1 ]] && echo "START DELETION:"
            deleting_password_reset_data=1
        fi
        if [[ "$line" =~ '^\s*$' ]]; then
            deleting_browser_session_data=0
            deleting_password_reset_data=0
        fi
        if [[ "$line" = '' ]]; then
            deleting_browser_session_data=0
            deleting_password_reset_data=0
        fi

        if [[ $deleting_browser_session_data = 1 ]]; then
            # do not write line to output
            [[ $verbose = 1 ]] && echo "DELETE $line"
        elif [[ $deleting_password_reset_data = 1 ]]; then
            # do not write line to output
            [[ $verbose = 1 ]] && echo "DELETE $line"
        else
            echo "$line" >> $tempfile
        fi
    done < "$input_file"

    # sanitize sensitive content
    # initialize tracking variables
    sanitize_expenditure=0
    sanitize_payment=0
    sanitize_session=0
    sanitize_user=0

    while IFS= read -r line
    do  
        # start sanitizing
        if [[ "$line" =~ 'INSERT INTO `expenditure`' ]]; then
            [[ $verbose = 1 ]] && echo "START SANITIZING:"
            sanitize_expenditure=1
        fi
        if [[ "$line" =~ 'INSERT INTO `payment`' ]]; then
            [[ $verbose = 1 ]] && echo "START SANITIZING:"
            sanitize_payment=1
        fi
        if [[ "$line" =~ 'INSERT INTO `session`' ]]; then
            [[ $verbose = 1 ]] && echo "START SANITIZING:"
            sanitize_session=1
        fi
        if [[ "$line" =~ 'INSERT INTO `user`' ]]; then
            [[ $verbose = 1 ]] && echo "START SANITIZING:"
            sanitize_user=1
        fi

        # stop sanitizing
        if [[ "$line" =~ '^\s*$' ]]; then
            sanitize_expenditure=0
            sanitize_payment=0
            sanitize_session=0
            sanitize_user=0
        fi
        if [[ "$line" = '' ]]; then
            sanitize_expenditure=0
            sanitize_payment=0
            sanitize_session=0
            sanitize_user=0
        fi

        # sanitize expenditures
        if [[ $sanitize_expenditure = 1 ]]; then
            # replace the comment column
            [[ $verbose = 1 ]] && echo "Replace line:"
            [[ $verbose = 1 ]] && echo $line
            [[ $verbose = 1 ]] && echo "By line:"
            [[ $verbose = 1 ]] && echo $line | perl -pe "s/,\s'[^']*'\)/, 'my comment'\)/"
            echo $line | perl -pe "s/,\s'[^']*'\)/, 'my comment'\)/" >> $output_file
        elif [[ $sanitize_payment = 1 ]]; then
            # replace the comment column
            [[ $verbose = 1 ]] && echo "Replace line:"
            [[ $verbose = 1 ]] && echo $line
            [[ $verbose = 1 ]] && echo "By line:"
            [[ $verbose = 1 ]] && echo $line | perl -pe "s/,\s'[^']*'\)/, 'my comment'\)/"
            echo $line | perl -pe "s/,\s'[^']*'\)/, 'my comment'\)/" >> $output_file
        elif [[ $sanitize_session = 1 ]]; then
            # replace the title and description column
            [[ $verbose = 1 ]] && echo "Replace line:"
            [[ $verbose = 1 ]] && echo $line
            [[ $verbose = 1 ]] && echo "By line:"
            [[ $verbose = 1 ]] && echo $line | perl -pe "s/,\s'[^']*',\s'[^']*',\s(\d+,\s\d+,\s\d+)/, 'my title', 'my description', \$1/"
            echo $line | perl -pe "s/,\s'[^']*',\s'[^']*',\s(\d+,\s\d+,\s\d+)/, 'my title', 'my description', \$1/" >> $output_file
        elif [[ $sanitize_user = 1 ]]; then
            # replace the title and description column
            [[ $verbose = 1 ]] && echo "Replace line:"
            [[ $verbose = 1 ]] && echo $line
            [[ $verbose = 1 ]] && echo "By line:"
            [[ $verbose = 1 ]] && echo $line | perl -pe "s/\((\d+),.*,(\s\d+,\s\d+,\s\d+,\s'[^']*')/\(\$1, 'username\$1', 12953, '\\\$6\\\$rounds=5000\\\$12953\\\$GmpDFyZnmBTriCjrIlT3tuWkWc5dvhHJeP76qt6ZWTQosQGeViJClsxvSlXesm9nw0Cs4R7BOQY5bqEqT.2nI.', NULL, 'firstname\$1', 'lastname\$1', 'street \$1', 'city \$1', \$1, 123456789\$1, 'email\$1\@domain\$1.com',\$2/"
            echo $line | perl -pe "s/\((\d+),.*,(\s\d+,\s\d+,\s\d+,\s'[^']*')/\(\$1, 'username\$1', 12953, '\\\$6\\\$rounds=5000\\\$12953\\\$GmpDFyZnmBTriCjrIlT3tuWkWc5dvhHJeP76qt6ZWTQosQGeViJClsxvSlXesm9nw0Cs4R7BOQY5bqEqT.2nI.', NULL, 'firstname\$1', 'lastname\$1', 'street \$1', 'city \$1', \$1, 123456789\$1, 'email\$1\@domain\$1.com',\$2/" >> $output_file
        else
            echo "$line" >> $output_file
        fi

    done < "$tempfile"

    rm $tempfile

    cat << 'EOF'

Anonymization of database dump has been done. Please manually check the dump as this script might not catch all user data.
EOF
}

main() {
    # check arguments
    if [ ! -f $input_file ]; then
        >&2 echo "Input file not found"
        exit 1
    fi

    # anonymize testdata
    anonymize

}

while getopts "hvi:o:" opt "$@"; do
    case "$opt" in
    h|help)
        show_help
        exit 0
        ;;
    v|verbose) verbose=1
        ;;
    i|in) input_file=$OPTARG
        ;;
    o|out) output_file=$OPTARG
    esac
done

main

exit 0