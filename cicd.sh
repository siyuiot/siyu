#!/bin/bash
#use: ./cicd.sh deploy-bs sms,lbs debug,d,c,k,e
# curl -XPOST "https://huangzhifeng:xxxx@bamboo.sharkgulf.cn/rest/api/latest/queue/SG-CICD?bamboo.variable.PROJECTS=deploy-bs&bamboo.variable.APPS=bikesvc&bamboo.variable.ENVS=d"
set -e

time=`date +'%Y%m%d%H%M'`
project=$1 #${bamboo.PROJECTS}
apps=$2 #${bamboo.APPS}
envs=$3 #${bamboo.ENVS}
dockerhub=docker.sharkgulf.cn

currGoVer=`go version | cut -d " "  -f3`
buildGoVer="go1.14.9"
if [ $currGoVer != $buildGoVer ];then
  echo "error: go version should be ${buildGoVer}"
  exit 1
fi

if [ -z $project ];then
    echo "please input deploy project"
    exit 1
fi
if [ -z $apps ];then
    echo "please input deploy apps"
    exit 1
fi
if [ -z $envs ];then
    echo "please input deploy envs"
    exit 1
fi
if [ $apps == SIYU ];then
    apps="user"
fi
for env in ${envs//,/ };do
    if [ $env = debug ];then
        echo "Your branch is up to date with origin/debug."
    elif [ $env = d ];then
        echo "Your branch is up to date with origin/dev."
    elif [ $env = c ];then
        git checkout master
    elif [ $env = k ];then
        dockerhub=registry.cn-hongkong.aliyuncs.com
        git checkout master
    elif [ $env = e ];then
        dockerhub=registry.cn-hongkong.aliyuncs.com
        git checkout master
    fi

    for app in ${apps//,/ };do
        buildApp=$app
        go env -w GOPROXY=https://goproxy.io,direct
        go env -w GOPRIVATE=*.sharkgulf.cn
        go env -w GOSUMDB=off
        # git config --global url."https://wizard:xxx@gogs.sharkgulf.cn".insteadof "https://gogs.sharkgulf.cn" #go get私有仓库认证问题
        echo $app && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./app/$buildApp/$buildApp ./app/$buildApp
        docker build -t $dockerhub/shark-bs/$app:$time ./app/$buildApp
        docker push $dockerhub/shark-bs/$app:$time
        # docker tag $dockerhub/shark-bs/$app:$time $dockerhub/shark-bs/$app:latest
        # docker push $dockerhub/shark-bs/$app:latest
        kubectl --context=$env cluster-info
        for resource in configmap service ingress deployment hap;do
            if [ $resource = deployment ];then
            sed "s/$app:latest/$app:$time/g" $project/$env/$resource/$app.yml | kubectl --context=$env apply -f -
            elif [ -f $project/$env/$resource/$app.yml ];then
            kubectl --context=$env apply -f $project/$env/$resource/$app.yml
            fi
        done
    done
done
