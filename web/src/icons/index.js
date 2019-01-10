import Vue from 'vue'

import Dashboard from './Dashboard.vue'
import Send from './Send.vue'
import Project from './Project.vue'
import User from './User.vue'
import Shell from './Shell.vue'
import Server from './Server.vue'
import Setting from './Setting.vue'
import CaretRight from './CaretRight.vue'
import Square from './Square.vue'
import MenuFold from './MenuFold.vue'
import MenuUnfold from './MenuUnfold.vue'
import CaretDown from './CaretDown.vue'
import Team from './Team.vue'
import Logout from './Logout.vue'
import Question from './Question.vue'
import Delete from './Delete.vue'
import Branch from './Branch.vue'
import Verify from './Verify.vue'

const components = [
    Dashboard,
    Send,
    Project,
    User,
    Shell,
    Server,
    Setting,
    CaretRight,
    Square,
    MenuFold,
    MenuUnfold,
    CaretDown,
    Team,
    Logout,
    Question,
    Delete,
    Branch,
    Verify,
]

const install = function(Vue) {
    components.forEach(component => {
        Vue.component(component.name, component);
    });
}

export default {
    install,
}
