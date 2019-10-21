#!/bin/bash

SV_MODE="update"
C_MODE="update"

echo "* Reading AWS variables"
AWS_REGION=us-west-1

echo "* Reading Screaming Vortex variables"
SV_PLATFORM=go-1.13.1
SV_KEYVALUE=screamingvortex
SV_APPNAME=screamingvortex
SV_DOMAIN="screamingvortex"
SV_PROFILE=EBKoronusAccess

echo "* Reading Calixis variables"
C_PLATFORM=python-3.6
C_KEYVALUE=screamingvortex
C_APPNAME=calixis
C_DOMAIN="sectorsofinequity"
C_PROFILE=EBKoronusAccess

while [ "$1" != "" ]; do
  case $1 in
    -r | --region ) shift
                    AWS_REGION=$1
                    ;;
    -m | --mode ) shift
                  SV_MODE=$1
                  C_MODE=$1
                  ;;
    --update )  SV_MODE=UPDATE
                C_MODE=UPDATE
                ;;
    --create )  SV_MODE=CREATE
                C_MODE=CREATE
                ;;
    --destroy ) SV_MODE=DESTROY
                C_MODE=DESTROY
                ;;
    --sv_profile )  shift
                    SV_PROFILE=$1
                    ;;
    --c_profile ) shift
                  C_PROFILE=$1
                  ;;
    --sv_domain ) shift
                  SV_DOMAIN=$1
                  ;;
    --c_domain )  shift
                  C_DOMAIN=$1
                  ;;
    --sv_name )   shift
                  SV_APPNAME=$1
                  ;;
    --c_name )    shift
                  C_APPNAME=$1
                  ;;
    --sv_key )    shift
                  SV_KEYVALUE=$1
                  ;;
    --c_key )     shift
                  C_KEYVALUE=$1
  esac
  shift
done

echo "* Goto Screaming Vortex"
cd ./screamingvortex
echo "* Initializing Screaming Vortex"
eb init -r $AWS_REGION -p $SV_PLATFORM -k $SV_KEYVALUE $SV_APPNAME
echo "* Building Screaming Vortex executable"
go build -o bin/application application.go
if [ ${SV_MODE^^} = "CREATE" ]
then
  echo "* Create Screaming Vortex"
  SV_PROFILE_ARN=$(aws iam get-instance-profile --instance-profile-name $SV_PROFILE | jq '.InstanceProfile.Arn')
  echo "=eb create $SV_APPNAME -ip $SV_PROFILE_ARN -c $SV_DOMAIN"
  eb create $SV_APPNAME -ip $SV_PROFILE -c $SV_DOMAIN
fi
if [ ${SV_MODE^^} = "CREATE" ] || [ ${SV_MODE^^} = "UPDATE" ]
then
  echo "* Deploy Screaming Vortex"
  eb deploy $SV_APPNAME
fi
if [ ${SV_MODE^^} = "DESTROY" ]
then
  echo "* Terminating Screaming Vortex"
  eb terminate $SV_APPNAME
fi

echo "* Goto Calixis"
cd ../calixis
echo "* Initializing Calixis"
eb init -r $AWS_REGION -p $C_PLATFORM -k $C_KEYVALUE $C_APPNAME
if [ ${C_MODE^^} = "CREATE" ]
then
  echo "* Create Calixis"
  C_PROFILE_ARN=$(aws iam get-instance-profile --instance-profile-name $C_PROFILE | jq '.InstanceProfile.Arn')
  eb create $C_APPNAME -ip $C_PROFILE_ARN -c $C_DOMAIN
  echo "* Set Calixis EnvVars"
  eb setenv SCREAMING_VORTEX_HOST=$SV_DOMAIN.$AWS_REGION.elasticbeanstalk.com
  echo "* Fetch Calixis CNAME"
  C_HOST="$(eb status $C_APP_NAME | grep CNAME: | awk '{print $2}')"
  echo "C_HOST=$C_HOST"
  echo "* Adding ALLOWED_HOSTS to Calixis"
  sed -i -e "s/^ALLOWED_HOSTS\s=\s\[.*]/ALLOWED_HOSTS = [\'localhost\', \'127.0.0.1\', \'$C_HOST\']/" ./calixis/settings.py
fi
if [ ${C_MODE} = "CREATE" ] || [ ${C_MODE} = "UPDATE" ]
then
  echo "* Deploy Calixis"
  eb deploy $C_APPNAME
fi
if [ ${C_MODE} = "DESTROY" ]
then
  echo "* Terminating Calixis"
  eb terminate $C_APPNAME
fi

echo "* Return to top directory"
cd ..
