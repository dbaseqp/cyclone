USER=$1
PASSWORD=$2
ADMINPW=$3

sed -e "s/{USER}/$USER/" \
    -e "s/{PASSWORD}/$PASSWORD/" \
    ./ldap/user_template.ldif > ./ldap/"$USER"_user_add.ldif
ldapadd -w $ADMINPW -H ldap://localhost:389 -D 'cn=admin,dc=kamino,dc=labs' -f ./ldap/"$USER"_user_add.ldif
rm ./ldap/"$USER"_user_add.ldif